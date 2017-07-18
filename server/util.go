package server

import (
	"net/http"
	"path"
	"strconv"

	"github.com/flosch/pongo2"
)

// render takes the provided context and passes it to the specified template
// for rendering.
func (s *Server) render(w http.ResponseWriter, templateName string, ctx pongo2.Context) {
	b, err := ReadFile(path.Join("templates", templateName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := pongo2.FromBytes(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	d, err := t.ExecuteBytes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(d)))
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(d)
}
