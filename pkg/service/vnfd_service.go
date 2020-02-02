package service

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/vishwanathj/protovnfdparser/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/vishwanathj/protovnfdparser/pkg/config"
	"github.com/vishwanathj/protovnfdparser/pkg/constants"
	"github.com/vishwanathj/protovnfdparser/pkg/dataaccess"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	"github.com/vishwanathj/protovnfdparser/pkg/mongo"
	"github.com/vishwanathj/protovnfdparser/pkg/utils"
)

type VnfdService struct {
	dal dataaccess.VnfdRepository
	cfg config.Config
}

var svcInstance *VnfdService
var once sync.Once

// GetVnfdServiceInstance returns a singleton instance of VnfdService
func GetVnfdServiceInstance(cfg config.Config) (*VnfdService, error) {
	once.Do(func() {
		d, err := mongo.NewMongoDAL(cfg)
		if err != nil {
			panic(err)
		}
		svcInstance = &VnfdService{d, cfg}
	})
	return svcInstance, nil
}

// CreateVnfd method that creates a new VNFD
func (p *VnfdService) CreateVnfd(u *models.Vnfd) errors.VnfdsvcError {
	log.Debug()

	// validate vnfd post body received from end user
	inputVnfd, err := json.Marshal(u)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusBadRequest}
	}
	log.Info(string(inputVnfd))

	err = utils.ValidateVnfdPostBody(inputVnfd)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusBadRequest}
	}

	u.SetCreationTimeAttributes()
	jsonval, err := json.Marshal(u)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"vnfdStrJsonVal": string(jsonval)}).Debug()

	//Before posting to the database, make sure, the newly created model adheres to the schema
	err = utils.ValidateVnfdInstanceBody(jsonval)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusBadRequest}
	}

	err = p.dal.InsertVnfd(u)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusConflict}
	}
	return errors.VnfdsvcError{nil, http.StatusOK}
}

// GetByVnfdname method that retrieves a VNFD given the name
func (p *VnfdService) GetByVnfdname(vnfdname string) (*models.Vnfd, error) {
	log.Debug()

	model, err := p.dal.FindVnfdByName(vnfdname)

	if err != nil {
		log.Error("Query Error:", err)
		return nil, err
	}
	jsonval, err := json.Marshal(model)
	if err != nil {
		log.Error("JSON Marshall error:", err)
		return nil, err
	}
	err = utils.ValidateVnfdInstanceBody(jsonval)
	if err != nil {
		log.Error("Failed ValidateVnfdInstanceBody:", err)
		return nil, err
	}
	log.WithFields(log.Fields{"vnfdname": vnfdname}).Debug()
	return model, err
}

// GetByVnfdID method that retrieves a VNFD given its ID
func (p *VnfdService) GetByVnfdID(vnfdID string) (*models.Vnfd, error) {
	log.Debug()
	model, err := p.dal.FindVnfdByID(vnfdID)
	if err != nil {
		log.Error("Query Error:", err)
		return nil, err
	}
	jsonval, err := json.Marshal(model)
	if err != nil {
		log.Error("JSON Marshall error:", err)
		return nil, err
	}
	err = utils.ValidateVnfdInstanceBody(jsonval)
	if err != nil {
		log.Error("Failed ValidateVnfdInstanceBody:", err)
		return nil, err
	}
	log.WithFields(log.Fields{"vnfdid": vnfdID}).Debug()
	return model, err
}

// GetVnfds method that retrieves a paginated list of Vnfds
func (p *VnfdService) GetVnfds(start string, limitinp int, sort string) (models.PaginatedVnfds, error) {
	log.Debug()

	var limit int

	if limitinp <= 0 || limitinp > p.cfg.PgntConfig.MaxLimit || limitinp < p.cfg.PgntConfig.MinLimit {
		limit = p.cfg.PgntConfig.DefaultLimit
		log.Debug("DefaultLimit being used instead of user provided input")
	} else {
		limit = limitinp
	}

	var res models.PaginatedVnfds
	vnfds, count, err := p.dal.GetVnfds(start, limit, sort)
	log.WithFields(log.Fields{"VNFDS": vnfds}).Debug("GET_VNFDS")

	if err != nil {
		return res, err
	}

	// Below validates each item in array adheres to schema and has not been manually tampered in DB directly
	// to be schema non-compliant
	for i := 0; i < len(vnfds); i++ {
		jsonval, err := json.Marshal(vnfds[i])
		if err != nil {
			log.Error("JSON Marshall error:", err)
			return res, err
		}
		err = utils.ValidateVnfdInstanceBody(jsonval)
		if err != nil {
			log.Error("Failed ValidateVnfdInstanceBody:", err)
			return res, err
		}
	}

	pgnCfg := p.cfg.PgntConfig
	var first *models.Link
	var next *models.Link
	var queryParams models.PgnQueryParams

	if len(sort) != 0 {
		queryParams = models.PgnQueryParams{OrderBy: sort}
	}

	first = &models.Link{Href: models.MakeFirstHref(*pgnCfg, limit, constants.ApiPathVnfds, queryParams)}

	if len(vnfds) == 0 || len(vnfds) < limit {
		next = nil
	} else {
		next = &models.Link{Href: models.MakeNextHref(*pgnCfg, limit, vnfds[limit-1].Name, constants.ApiPathVnfds, queryParams)}
	}
	res = models.PaginatedVnfds{Limit: limit, TotalCount: count, First: first, Next: next, Vnfds: vnfds}

	// Ensure that the resulting container in "res" validates against the defined schema
	//jsonval, err := json.Marshal(vnfds)
	jsonval, err := json.Marshal(res)
	if err != nil {
		log.Error("JSON Marshall error:", err)
		return res, err
	}
	err = utils.ValidatePaginatedVnfdsInstancesBody(jsonval)
	return res, err
}

// GetInputParamsSchemaForVnfd method that returns valid InputParams schema given a Vnfd
func (p *VnfdService) GetInputParamsSchemaForVnfd(jsonval []byte) ([]byte, error) {
	log.Debug()
	inp, err := utils.GenerateJSONSchemaFromParameterizedTemplate(jsonval)
	log.WithFields(log.Fields{"inp": string(inp)}).Debug("Inputs received from end user")
	return inp, err
}

// GetHealth method used for liveness probe by kubernetes
func (p *VnfdService) GetHealth() string {
	log.Debug()
	return "OK"
}

// GetReadiness method use for mongodb readiness probe by kubernetes
/*func (p *VnfdService) GetReadiness() (string) {
	log.Debug()
	p.collection.Find(bson.M{"name": bson.M{"$gt": ""}})
	_, err := p.collection.Count()
	if err != nil {
		return "FAIL"
	}
	return "OK"
}*/
