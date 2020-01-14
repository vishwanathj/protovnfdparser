package server

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/vishwanathj/protovnfdparser/pkg/dataaccess"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	utils "github.com/vishwanathj/protovnfdparser/pkg/utils"
	//utils "github.com/vishwanathj/JSON-Parameterized-Data-Validator/pkg/jsondatavalidator"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v2"
)

// VnfdIDPattern regexp pattern for the VnfdId
const VnfdIDPattern = "^VNFD-[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"

type vnfdRouter struct {
	vnfdService models.VnfdService
	dalService  dataaccess.DataAccessLayer
}

// NewVnfdRouter creates a new mux router
func NewVnfdRouter(v models.VnfdService, d dataaccess.DataAccessLayer, router *mux.Router) *mux.Router {
	log.Debug()
	vnfdRouter := vnfdRouter{v, d}

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

func (ur *vnfdRouter) createVnfdHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Make a clone for later processing
	bodyClone := ioutil.NopCloser(bytes.NewBuffer(body))

	log.WithFields(log.Fields{"POST_BODY_VNFD_ORIG": string(body)}).Debug()
	err = utils.ValidateVnfdPostBody(body)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	//restore the body for the http.Request
	r.Body = bodyClone
	vnfd, err := decodeVnfd(r)
	log.WithFields(log.Fields{"Vnfd": vnfd}).Debug()
	if err != nil {
		log.WithFields(log.Fields{"decodeVnfdErr": err}).Debug()
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = ur.vnfdService.CreateVnfd(&vnfd)
	if err != nil {
		log.WithFields(log.Fields{"CreateVnfdErr": err}).Error()
		Error(w, http.StatusConflict, err.Error())
		return
	}

	JSON(w, http.StatusOK, err)
}

func (ur *vnfdRouter) getVnfdsHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	vars := r.URL.Query()
	start := vars.Get("start")
	limit := vars.Get("limit")

	log.WithFields(log.Fields{"LIMIT": limit, "START": start}).Debug("Inputs received from end user")

	l, err := strconv.Atoi(limit)
	if err != nil || l <= 0 || l > models.MaxLimit || l < models.MinLimit {
		l = models.DefaultLimit
		log.Debug("DefaultLimit being used instead of user provided input")
	}

	vnfds, err := ur.vnfdService.GetVnfds(start, l)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, vnfds)
}

func (ur *vnfdRouter) getVnfdHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	vars := mux.Vars(r)
	vnfdname := vars["name"]

	validVnfdID := regexp.MustCompile(VnfdIDPattern)
	if validVnfdID.MatchString(vnfdname) {
		vnfd, err := ur.vnfdService.GetByVnfdID(vnfdname)
		if err != nil {
			Error(w, http.StatusNotFound, err.Error())
			return
		}
		JSON(w, http.StatusOK, vnfd)
	} else {
		vnfd, err := ur.vnfdService.GetByVnfdname(vnfdname)
		if err != nil {
			Error(w, http.StatusNotFound, err.Error())
			return
		}
		JSON(w, http.StatusOK, vnfd)
	}
}

func (ur *vnfdRouter) getVnfdInputParamsSchemaHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	vars := mux.Vars(r)
	vnfdname := vars["name"]

	var vnfd *models.Vnfd
	var err error
	var jsonval []byte
	var inputparam []byte

	validVnfdID := regexp.MustCompile(VnfdIDPattern)
	if validVnfdID.MatchString(vnfdname) {
		vnfd, err = ur.vnfdService.GetByVnfdID(vnfdname)
	} else {
		vnfd, err = ur.vnfdService.GetByVnfdname(vnfdname)
	}

	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	jsonval, err = yaml.Marshal(vnfd)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	inputparam, err = ur.vnfdService.GetInputParamsSchemaForVnfd(jsonval)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	var m map[string]interface{}
	err = json.Unmarshal(inputparam, &m)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Debug(m)
	res := m["properties"].(map[string]interface{})
	res["vnfd_id"] = vnfd.ID
	log.Debug(res)
	JSON(w, http.StatusOK, res)
}

func (ur *vnfdRouter) livenessProbe(w http.ResponseWriter, r *http.Request) {
	log.Debug()
	JSON(w, http.StatusOK, ur.vnfdService.GetHealth())
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
