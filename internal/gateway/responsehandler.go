package gateway

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	response := make(map[string]interface{})

	if status < 400 {
		// positive response
		response["status"] = true
		response["data"] = payload
	} else {
		// error response
		response["status"] = false
		response["errors"] = payload
	}

	jsonData, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(jsonData))
}
