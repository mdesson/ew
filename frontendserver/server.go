package frontendserver

import (
	"encoding/json"
	"fmt"
	"github.com/mdesson/ew/consts"
	"github.com/mdesson/ew/echoserver"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
)

type Server struct {
	l        *slog.Logger
	port     int
	reqChan  chan echoserver.RequestDetails
	Requests []echoserver.RequestDetails
}

func New(port int, l *slog.Logger, reqChan chan echoserver.RequestDetails) *Server {
	return &Server{port: port, l: l, Requests: make([]echoserver.RequestDetails, 0), reqChan: reqChan}
}

func (s *Server) Start() {
	// send request details to the web UI server
	go func() {
		// TODO: add .Stop() function which closes the channel
		for {
			select {
			case r := <-s.reqChan:
				s.Requests = append(s.Requests, r)
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/requests", s.getRequestsHandler)
	mux.HandleFunc("/", serveIndexHTML)

	s.l.With(consts.FieldFunction, "Start")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux))
}

func (s *Server) getRequestsHandler(w http.ResponseWriter, r *http.Request) {
	l := s.l.With(consts.FieldFunction, "getRequestsHandler")
	if r.Method != http.MethodGet {
		l.Error("method should be GET", "error", fmt.Errorf("Client sent %s %s", r.Method, r.URL.String()))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := json.Marshal(s.Requests)
	if err != nil {
		l.Error("failed to marshal requests", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err := w.Write(body); err != nil {
		l.Error("failed to write response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join("public", "index.html")
	http.ServeFile(w, r, filePath)
}
