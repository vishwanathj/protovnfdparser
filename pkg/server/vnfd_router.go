package server

import (
	"encoding/json"
	"errors"

	"github.com/vishwanathj/protovnfdparser/pkg/constants"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	vsvcerr "github.com/vishwanathj/protovnfdparser/pkg/errors"
	"github.com/vishwanathj/protovnfdparser/pkg/models"

	"net/http"
	//"regexp"
	"strconv"

	"gopkg.in/yaml.v2"
)

type vnfdRouter struct {
	vnfdService models.VnfdService
}

// NewVnfdRouter creates a new mux router
func NewVnfdRouter(v models.VnfdService, router *mux.Router) *mux.Router {
	log.Debug()
	vnfdRouter := vnfdRouter{v}

	/*router.HandleFunc("", vnfdRouter.createVnfdHandler).Methods("POST")
	router.HandleFunc("", vnfdRouter.getVnfdsHandler).Methods("GET")
	router.HandleFunc("/{name}", vnfdRouter.getVnfdHandler).Methods("GET")
	router.HandleFunc("/{name}/inputparamschema", vnfdRouter.getVnfdInputParamsSchemaHandler).Methods("GET")
	router.HandleFunc("/health", vnfdRouter.livenessProbe).Methods("GET")*/
	router.HandleFunc("/vnfds", vnfdRouter.createVnfdHandler).Methods("POST")
	router.HandleFunc("/vnfds", vnfdRouter.getVnfdsHandler).Methods("GET")
	router.HandleFunc("/vnfds/{name}", vnfdRouter.getVnfdHandler).Methods("GET")
	router.HandleFunc("/vnfds/{name}/inputparamschema", vnfdRouter.getVnfdInputParamsSchemaHandler).Methods("GET")
	router.HandleFunc("/health", vnfdRouter.livenessProbe).Methods("GET")
	//router.HandleFunc("/readiness", vnfdRouter.readinessProbe).Methods("GET")
	return router
}

func (vr *vnfdRouter) createVnfdHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	vnfd, err := decodeVnfd(r)
	log.WithFields(log.Fields{"Vnfd": vnfd}).Debug()
	if err != nil {
		log.WithFields(log.Fields{"decodeVnfdErr": err}).Debug()
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	svcerr := vr.vnfdService.CreateVnfd(&vnfd)
	if svcerr.OrigError != nil {
		log.WithFields(log.Fields{"CreateVnfdErr": err}).Error()
		Error(w, int(svcerr.HttpCode), svcerr.OrigError.Error())
		return
	}

	JSON(w, http.StatusOK, err)
}

func (vr *vnfdRouter) getVnfdsHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	vars := r.URL.Query()
	//start := vars.Get("start")
	//limit := vars.Get("limit")
	start := vars.Get(constants.PaginationURLStart)
	limit := vars.Get(constants.PaginationURLLimit)
	sort := vars.Get(constants.PaginationURLSort)

	log.WithFields(log.Fields{"LIMIT": limit, "START": start, "SORT": sort}).Debug("Inputs received from end user")

	var l int
	var err error
	if len(limit) == 0 {
		l = 0
	} else {
		l, err = strconv.Atoi(limit)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	var svcerr vsvcerr.VnfdsvcError
	vnfds, svcerr := vr.vnfdService.GetVnfds(start, l, sort)
	if svcerr.OrigError != nil {
		Error(w, int(svcerr.HttpCode), err.Error())
		return
	}

	JSON(w, http.StatusOK, vnfds)
}

func (vr *vnfdRouter) getVnfdHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	vars := mux.Vars(r)
	vnfdname := vars["name"]

	vnfd, err := vr.vnfdService.GetVnfd(vnfdname)
	if err.OrigError != nil {
		Error(w, int(err.HttpCode), err.Error())
		return
	}
	JSON(w, http.StatusOK, vnfd)

	/*validVnfdID := regexp.MustCompile(constants.VnfdIDPattern)
	if validVnfdID.MatchString(vnfdname) {
		vnfd, err := vr.vnfdService.GetByVnfdID(vnfdname)
		if err.OrigError != nil {
			Error(w, int(err.HttpCode), err.Error())
			return
		}
		JSON(w, http.StatusOK, vnfd)
	} else {
		vnfd, err := vr.vnfdService.GetByVnfdname(vnfdname)
		if err.OrigError != nil {
			Error(w, int(err.HttpCode), err.Error())
			return
		}
		JSON(w, http.StatusOK, vnfd)
	}*/
}

func (vr *vnfdRouter) getVnfdInputParamsSchemaHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	vars := mux.Vars(r)
	vnfdname := vars["name"]

	var vnfd *models.Vnfd
	var err vsvcerr.VnfdsvcError
	var jsonval []byte
	var inputparam []byte

	vnfd, err = vr.vnfdService.GetVnfd(vnfdname)
	/*validVnfdID := regexp.MustCompile(constants.VnfdIDPattern)
	if validVnfdID.MatchString(vnfdname) {
		vnfd, err = vr.vnfdService.GetByVnfdID(vnfdname)
	} else {
		vnfd, err = vr.vnfdService.GetByVnfdname(vnfdname)
	}*/

	if err.OrigError != nil {
		Error(w, int(err.HttpCode), err.Error())
		return
	}
	jsonval, errg := yaml.Marshal(vnfd)
	if errg != nil {
		Error(w, http.StatusInternalServerError, errg.Error())
		return
	}
	inputparam, svcerrg := vr.vnfdService.GetInputParamsSchemaForVnfd(jsonval)
	if svcerrg.OrigError != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	var m map[string]interface{}
	errg = json.Unmarshal(inputparam, &m)
	if errg != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Debug(m)
	res := m["properties"].(map[string]interface{})
	res["vnfd_id"] = vnfd.ID
	log.Debug(res)
	JSON(w, http.StatusOK, res)
}

func (vr *vnfdRouter) livenessProbe(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	JSON(w, http.StatusOK, vr.vnfdService.GetHealth())
}

/*func (ur *vnfdRouter) readinessProbe(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	JSON(w, http.StatusOK, ur.vnfdService.GetReadiness())
}*/

func decodeVnfd(r *http.Request) (models.Vnfd, error) {
	log.Debug()
	var v models.Vnfd
	if r.Body == nil {
		return v, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)
	return v, err
}
