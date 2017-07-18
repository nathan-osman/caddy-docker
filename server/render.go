package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/flosch/pongo2"
)

// b0xLoader provides a Pongo2 loader for b0x template files.
type b0xLoader struct{}

// Abs returns the absolute path to a template file.
func (l b0xLoader) Abs(base, name string) string {
	return name
}

// Get retrieves a reader for the specified path.
func (l b0xLoader) Get(path string) (io.Reader, error) {
	f, err := FS.OpenFile(
		CTX,
		fmt.Sprintf("templates/%s", path),
		os.O_RDONLY,
		0,
	)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

var (
	templateLoader = b0xLoader{}
	templateSet    = pongo2.NewSet("", templateLoader)
)

// render takes the provided context and passes it to the specified template
// for rendering.
func (s *Server) render(w http.ResponseWriter, templateName string, ctx pongo2.Context) {
	t, err := templateSet.FromFile(templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := t.ExecuteBytes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
