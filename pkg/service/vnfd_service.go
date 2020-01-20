package service

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/vishwanathj/protovnfdparser/pkg/config"
	"github.com/vishwanathj/protovnfdparser/pkg/constants"
	"github.com/vishwanathj/protovnfdparser/pkg/dataaccess"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	"github.com/vishwanathj/protovnfdparser/pkg/mongo"
	"github.com/vishwanathj/protovnfdparser/pkg/utils"
)

type VnfdSvc struct {
	dal dataaccess.VnfdRepository
	cfg config.Config
}

var svcInstance *VnfdSvc
var once sync.Once

// GetVnfdServiceInstance returns a singleton instance of VnfdSvc
func GetVnfdServiceInstance(cfg config.Config) (*VnfdSvc, error) {
	once.Do(func() {
		d, err := mongo.NewMongoDAL(cfg)
		if err != nil {
			panic(err)
		}
		svcInstance = &VnfdSvc{d, cfg}
	})
	return svcInstance, nil
}

// CreateVnfd method that creates a new VNFD
func (p *VnfdSvc) CreateVnfd(u *models.Vnfd) error {
	log.Debug()

	u.SetCreationTimeAttributes()
	jsonval, err := json.Marshal(u)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"vnfdStrJsonVal": string(jsonval)}).Debug()

	//Before posting to the database, make sure, the newly created model adheres to the schema
	err = utils.ValidateVnfdInstanceBody(jsonval)
	if err != nil {
		//log.WithFields(log.Fields{"CreateVnfdError": err}).Error()
		return err
	}

	return p.dal.InsertVnfd(u)
}

// GetByVnfdname method that retrieves a VNFD given the name
func (p *VnfdSvc) GetByVnfdname(vnfdname string) (*models.Vnfd, error) {
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
func (p *VnfdSvc) GetByVnfdID(vnfdID string) (*models.Vnfd, error) {
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
func (p *VnfdSvc) GetVnfds(start string, limitinp int) (models.PaginatedVnfds, error) {
	log.Debug()

	var limit int

	if limitinp <= 0 || limitinp > p.cfg.PgntConfig.MaxLimit || limitinp < p.cfg.PgntConfig.MinLimit {
		limit = p.cfg.PgntConfig.DefaultLimit
		log.Debug("DefaultLimit being used instead of user provided input")
	} else {
		limit = limitinp
	}

	var res models.PaginatedVnfds
	vnfds, count, err := p.dal.GetVnfds(start, limit)
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
	if len(vnfds) == 0 || len(vnfds) < limit {
		first := models.Link{Href: models.MakeFirstHref(*pgnCfg, limit, constants.ApiPathVnfds)}
		res = models.PaginatedVnfds{Limit: limit, TotalCount: count, First: &first, Next: nil, Vnfds: vnfds}
	} else {
		//log.WithFields(log.Fields{"LAST": vnfds[limit-1].Name})
		first := models.Link{Href: models.MakeFirstHref(*pgnCfg, limit, constants.ApiPathVnfds)}
		next := models.Link{Href: models.MakeNextHref(*pgnCfg, limit, vnfds[limit-1].Name, constants.ApiPathVnfds)}
		res = models.PaginatedVnfds{Limit: limit, TotalCount: count, First: &first, Next: &next, Vnfds: vnfds}
	}

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
func (p *VnfdSvc) GetInputParamsSchemaForVnfd(jsonval []byte) ([]byte, error) {
	log.Debug()
	inp, err := utils.GenerateJSONSchemaFromParameterizedTemplate(jsonval)
	log.WithFields(log.Fields{"inp": string(inp)}).Debug("Inputs received from end user")
	return inp, err
}

// GetHealth method used for liveness probe by kubernetes
func (p *VnfdSvc) GetHealth() string {
	log.Debug()
	return "OK"
}

// GetReadiness method use for mongodb readiness probe by kubernetes
/*func (p *VnfdSvc) GetReadiness() (string) {
	log.Debug()
	p.collection.Find(bson.M{"name": bson.M{"$gt": ""}})
	_, err := p.collection.Count()
	if err != nil {
		return "FAIL"
	}
	return "OK"
}*/
