package echoserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdesson/ew/consts"
	"github.com/mdesson/ew/util"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	l       *slog.Logger
	port    int
	reqChan chan RequestDetails
}

func New(port int, l *slog.Logger, reqChan chan RequestDetails) *Server {
	return &Server{port: port, l: l, reqChan: reqChan}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.listenHandler)

	s.l.With(consts.FieldFunction, "start")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux))
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
		s.reqChan <- details
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
		fmt.Println("BOOOOOOOOOP")
		err = errors.New("multipart/form-data is not supported")
	default:
		str, err = printBodyToString(r)
	}

	if err != nil {
		l.Error(err.Error())
		details.ErrorMsg = util.Ptr(err.Error())
	} else if str != "" {
		details.Body = &str
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
