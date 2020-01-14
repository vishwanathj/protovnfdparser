package server

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Error ...
func Error(w http.ResponseWriter, code int, message string) {
	log.Debug()
	JSON(w, code, map[string]string{"error": message})
}

// JSON ...
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	log.Debug()
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	log.WithFields(log.Fields{"response": payload}).Debug("JSON: Write Response")
}
