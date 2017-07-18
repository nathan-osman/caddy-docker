package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server provides a web interface for interacting with the application.
type Server struct {
	username string
	password string
	router   *mux.Router
	listener net.Listener
	log      *logrus.Entry
	stopped  chan bool
}

// New creates a new server.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		srv    = http.Server{}
		router = mux.NewRouter()
		s      = &Server{
			username: cfg.Username,
			password: cfg.Password,
			router:   router,
			listener: l,
			log:      logrus.WithField("context", "server"),
			stopped:  make(chan bool),
		}
	)
	srv.Handler = s
	router.HandleFunc("/", s.index)
	router.PathPrefix("/static").Handler(http.FileServer(HTTP))
	go func() {
		defer close(s.stopped)
		defer s.log.Info("HTTP server has stopped")
		s.log.Info("starting HTTP server")
		if err := srv.Serve(l); err != nil {
			s.log.Error(err)
		}
	}()
	return s, nil
}

// ServeHTTP ensures that HTTP basic auth credentials match the requires ones
// and then routes the request to the muxer.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(s.username) != 0 && len(s.password) != 0 {
		u, p, ok := r.BasicAuth()
		if !ok || u != s.username || p != s.password {
			w.Header().Set("WWW-Authenticate", "Basic realm=caddy")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}
	s.router.ServeHTTP(w, r)
}

// Close shuts down the server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stopped
}
