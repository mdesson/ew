package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdesson/ew/consts"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	l        *slog.Logger
	port     int
	Requests []RequestDetails
}

func New(port int, l *slog.Logger) *Server {
	return &Server{port: port, l: l, Requests: make([]RequestDetails, 0)}
}

func (s *Server) Start() {
	http.HandleFunc("/", s.listenHandler)

	s.l.With(consts.FieldFunction, "start")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil))
}

func (s *Server) listenHandler(w http.ResponseWriter, r *http.Request) {
	l := s.l.With(consts.FieldFunction, "listenHandler")

	details := RequestDetails{
		Date:        time.Now(),
		Method:      r.Method,
		Path:        r.URL.Path,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}

	defer func() {
		s.Requests = append(s.Requests, details)
	}()

	// add headers
	for k, v := range r.Header {
		details.Headers[k] = strings.Join(v, " ")
	}

	// add query params
	for k, v := range r.URL.Query() {
		details.QueryParams[k] = strings.Join(v, " ")
	}

	// get the body and error, if any
	var str string
	var err error
	switch contentType := r.Header.Get("Content-Type"); contentType {
	case "application/json":
		str, err = printJSON(r)
	case "application/x-www-form-urlencoded":
		str, err = printURLEncoded(r)
	case "multipart/form-data":
		err = errors.New("multipart/form-data is not supported")
	default:
		str, err = printBodyToString(r)
	}

	if err != nil {
		l.Error(err.Error())
		details.Err = err
	} else {
		details.Body = str
	}

	l.Debug(details.String())

	w.WriteHeader(http.StatusOK)
}

func printJSON(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", err
	}

	if !json.Valid(body) {
		return "", fmt.Errorf("invalid json: %s", string(body))
	}

	prettyJSON := new(bytes.Buffer)
	if err := json.Indent(prettyJSON, body, "", "  "); err != nil {
		return "", nil
	}
	return fmt.Sprint(prettyJSON.String(), ""), nil
}

func printURLEncoded(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", nil
	}

	formValues := make([]string, 0)
	for key := range r.PostForm {
		formValues = append(formValues, fmt.Sprintf("%s=%s", key, r.PostFormValue(key)))
	}

	return strings.Join(formValues, "\n"), nil
}

func printBodyToString(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}
