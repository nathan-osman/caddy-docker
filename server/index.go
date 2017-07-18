package server

import (
	"net/http"

	"github.com/flosch/pongo2"
)

// index displays the home page.
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	s.render(w, "index.html", pongo2.Context{})
}
