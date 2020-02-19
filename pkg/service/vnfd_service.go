package service

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sync"

	"github.com/vishwanathj/protovnfdparser/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/vishwanathj/protovnfdparser/pkg/config"
	"github.com/vishwanathj/protovnfdparser/pkg/constants"
	"github.com/vishwanathj/protovnfdparser/pkg/datarepo"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	"github.com/vishwanathj/protovnfdparser/pkg/mongo"
	"github.com/vishwanathj/protovnfdparser/pkg/utils"
)

type VnfdService struct {
	dal datarepo.VnfdRepository
	cfg config.Config
}

var svcInstance *VnfdService
var once sync.Once
var vnfdVerifier = utils.NewVnfdValidator()

// GetVnfdServiceInstance returns a singleton instance of VnfdService
func GetVnfdServiceInstance(cfg config.Config) (*VnfdService, error) {
	once.Do(func() {
		d, err := mongo.NewMongoRepo(cfg)
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

	// validate vnfd post body received from end user; need to marshal to json first
	inputVnfd, err := json.Marshal(u)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusBadRequest}
	}
	log.Info(string(inputVnfd))

	err = vnfdVerifier.ValidateVnfdPostBody(inputVnfd)
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
	err = vnfdVerifier.ValidateVnfdInstanceBody(jsonval)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusBadRequest}
	}

	err = p.dal.InsertVnfd(u)
	if err != nil {
		return errors.VnfdsvcError{err, http.StatusConflict}
	}
	return errors.VnfdsvcError{nil, http.StatusOK}
}

// GetVnfd method that retrieves a VNFD given the name or the ID
func (p *VnfdService) GetVnfd(nameorid string) (*models.Vnfd, errors.VnfdsvcError) {
	log.Debug()

	var model *models.Vnfd
	var err error
	validVnfdID := regexp.MustCompile(constants.VnfdIDPattern)
	if validVnfdID.MatchString(nameorid) {
		model, err = p.dal.FindVnfdByID(nameorid)

	} else {
		model, err = p.dal.FindVnfdByName(nameorid)
	}
	if err != nil {
		log.Error("Query Error:", err)
		return nil, errors.VnfdsvcError{err, http.StatusNotFound}
	}
	jsonval, err := json.Marshal(model)
	if err != nil {
		log.Error("JSON Marshall error:", err)
		return nil, errors.VnfdsvcError{err, http.StatusNotFound}
	}
	err = vnfdVerifier.ValidateVnfdInstanceBody(jsonval)
	if err != nil {
		log.Error("Failed ValidateVnfdInstanceBody:", err)
		return nil, errors.VnfdsvcError{err, http.StatusInternalServerError}
	}
	log.WithFields(log.Fields{"nameorid": nameorid}).Debug()
	return model, errors.VnfdsvcError{nil, http.StatusOK}
}

// GetVnfds method that retrieves a paginated list of Vnfds
func (p *VnfdService) GetVnfds(start string, limitinp int, sort string) (models.PaginatedVnfds, errors.VnfdsvcError) {
	log.Debug()

	var limit int

	if limitinp <= 0 || limitinp > p.cfg.PgntConfig.MaxLimit || limitinp < p.cfg.PgntConfig.MinLimit {
		limit = p.cfg.PgntConfig.DefaultLimit
		log.Info("DefaultLimit being used instead of user provided input")
	} else {
		limit = limitinp
	}

	var res models.PaginatedVnfds
	vnfds, count, err := p.dal.GetVnfds(start, limit, sort)
	log.WithFields(log.Fields{"VNFDS": vnfds, "count": count, "err": err}).Info("GET_VNFDS")

	if err != nil {
		log.Info(err)
		return res, errors.VnfdsvcError{err, http.StatusInternalServerError}
	}

	// Below validates each item in array adheres to schema and has not been manually tampered in DB directly
	// to be schema non-compliant
	for i := 0; i < len(vnfds); i++ {
		jsonval, err := json.Marshal(vnfds[i])
		if err != nil {
			log.Info("JSON Marshall error:", err)
			return res, errors.VnfdsvcError{err, http.StatusInternalServerError}
		}
		err = vnfdVerifier.ValidateVnfdInstanceBody(jsonval)
		if err != nil {
			log.Info("Failed ValidateVnfdInstanceBody:", err)
			return res, errors.VnfdsvcError{err, http.StatusInternalServerError}
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
		log.Info("JSON Marshall error:", err)
		return res, errors.VnfdsvcError{err, http.StatusInternalServerError}
	}
	err = vnfdVerifier.ValidatePaginatedVnfdsInstancesBody(jsonval)
	if err != nil {
		return res, errors.VnfdsvcError{err, http.StatusInternalServerError}
	}
	return res, errors.VnfdsvcError{nil, http.StatusOK}
}

// GetInputParamsSchemaForVnfd method that returns valid InputParams schema given a Vnfd
func (p *VnfdService) GetInputParamsSchemaForVnfd(jsonval []byte) ([]byte, errors.VnfdsvcError) {
	log.Debug()
	inp, err := utils.GenerateJSONSchemaFromParameterizedTemplate(jsonval)
	log.WithFields(log.Fields{"inp": string(inp)}).Debug("Inputs received from end user")
	if err != nil {
		return nil, errors.VnfdsvcError{err, http.StatusInternalServerError}
	}
	return inp, errors.VnfdsvcError{nil, http.StatusOK}
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
