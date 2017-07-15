package configurator

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/mholt/caddy"
	"github.com/nathan-osman/caddy-docker/container"
)

var tmpl *template.Template

func init() {
	tmpl = template.New("caddy")
	tmpl.Funcs(map[string]interface{}{
		"join": func(list []string) string {
			return strings.Join(list, ", ")
		},
	})
	template.Must(tmpl.Parse(
		`# Auto generated file
{{range $c := .}}
# {{$c.Name}}
{{$c.Domains|join}} {
    proxy / {{$c.Addr}}
}
{{end}}
`,
	))
}

// generate enumerates all of the containers and builds a configuration file.
func (c *Configurator) generate(m map[string]*container.Container) error {
	c.log.Info("generating new configuration")
	w := bytes.NewBuffer(nil)
	if err := tmpl.ExecuteTemplate(w, "caddy", m); err != nil {
		return err
	}
	cdyfile := &caddy.CaddyfileInput{
		Filepath:       "internal",
		Contents:       w.Bytes(),
		ServerTypeName: "http",
	}
	return func() error {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		var (
			inst *caddy.Instance
			err  error
		)
		if c.inst != nil {
			c.log.Info("restarting server")
			inst, err = c.inst.Restart(cdyfile)
		} else {
			c.log.Info("starting server")
			inst, err = caddy.Start(cdyfile)
		}
		if err != nil {
			return err
		}
		c.inst = inst
		return nil
	}()
}
