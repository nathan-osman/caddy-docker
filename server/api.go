package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// makeError creates a map with a single key describing an error condition.
func makeError(desc string) interface{} {
	return map[string]interface{}{
		"error": desc,
	}
}

// api provides methods that the UI can utilize.
func (s *Server) api(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var response interface{}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Form.Get("action") {
	case "getContainers":
		response = s.configurator.Containers()
	case "restartContainer":
		id := r.Form.Get("id")
		if len(id) == 0 {
			response = makeError("container ID is missing")
		} else if err := s.monitor.Restart(r.Context(), id); err != nil {
			response = makeError(err.Error())
		} else {
			response = map[string]interface{}{}
		}
	default:
		response = makeError("invalid action")
	}
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
