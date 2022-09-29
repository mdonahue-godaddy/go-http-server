package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mdonahue-godaddy/go-http-server/config"
	"github.com/mdonahue-godaddy/go-http-server/metrics/gometrics"
	"github.com/mdonahue-godaddy/go-http-server/shared"
)

const (
	DefaultResponseTemplateFile string = "<!DOCTYPE HTML><html lang='en-us'><head><title>** {{page_title}} **</title></head><body>{{page_body}}</body></html>"
)

// Server struct
type Server struct {
	serviceName          string
	isInitialized        bool
	isShuttingDown       bool
	config               *config.Settings
	context              context.Context
	router               *http.ServeMux
	server               *http.Server
	metrics              gometrics.IGoMetrics
	responseTemplateFile string
}

// NewServer - create new instance of server
func NewServer(serviceName string, cfg *config.Settings, template *string) *Server {
	server := Server{}

	server.serviceName = serviceName
	server.config = cfg

	if template != nil && len(*template) > 0 {
		server.responseTemplateFile = *template
	} else {
		server.responseTemplateFile = DefaultResponseTemplateFile
	}

	server.metrics = gometrics.NewGoMetrics()

	return &server
}

func (s *Server) GetConfig() *config.Settings {
	return s.config
}

func (s *Server) GetServiceName() string {
	return s.serviceName
}

func (s *Server) IsInitialized() bool {
	return s.isInitialized
}

func (s *Server) IsShuttingDown() bool {
	return s.isShuttingDown
}

// doHeadErrorResponse - WARNING HEAD requests should not return a body so normal error response can't be used.
func (s *Server) DoHeadErrorResponse(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, httpStatusCode int, message string) {
	method := "server.doHeadRequestErrorResponse"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseStatusCode, httpStatusCode)).Debugf("%s entering", method)

	responseWriter.Header().Add("X-Error", fmt.Sprintf("HEAD not supported, %s", message))
	responseWriter.Header().Add("Content-Length", "0")
	responseWriter.Header().Add("Content-Type", shared.ContentType_TextHtml)

	s.WriteHeader(ctx, responseWriter, httpStatusCode)
}

// doErrorResponse - Error response processor
func (s *Server) DoErrorResponse(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, httpStatusCode int, htmlMessage string, err error) {
	method := "server.doErrorResponse"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseStatusCode, httpStatusCode, shared.KeyHTTPResponseBodyContent, htmlMessage)).Debugf("%s entering", method)

	shared.AddUniversalHeaders(ctx, responseWriter, s.serviceName)

	msg := ""
	if err != nil {
		msg = err.Error()
	}

	log.WithFields(shared.GetFields(ctx, shared.EventTypeError, false, shared.KeyHTTPResponseStatusCode, httpStatusCode, shared.KeyHTTPResponseBodyContent, htmlMessage)).Errorf("%s returning error response. %s", method, msg)

	s.WriteHeader(ctx, responseWriter, httpStatusCode)

	if len(htmlMessage) > 0 {
		err := shared.WriteHTML(ctx, responseWriter, htmlMessage)
		if err != nil {
			log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseStatusCode, httpStatusCode, shared.KeyHTTPResponseBodyContent, htmlMessage, shared.KeyErrorMessage, err.Error())).Errorf("%s error calling writeHTML", method)
		}
	}
}

func (s *Server) DoValidRequestResponse(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, responseHTML string) {
	method := "server.doValidRequestResponse"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s entering", method)

	shared.AddUniversalHeaders(ctx, responseWriter, s.serviceName)

	// Http Status Code 200
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s Masking", method)

	s.WriteHeader(ctx, responseWriter, http.StatusOK)

	err := shared.WriteHTML(ctx, responseWriter, responseHTML)
	if err != nil {
		log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyErrorMessage, err.Error())).Errorf("%s writeHTML error", method)
	}
}

func (s *Server) WriteHeader(ctx context.Context, responseWriter http.ResponseWriter, httpStatusCode int) {
	method := "server.writeHeader"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s writing response header with HTTP Status Code: %d", method, httpStatusCode)

	s.metrics.IncHTTPStatusCounters(ctx, httpStatusCode)

	responseWriter.WriteHeader(httpStatusCode)
}

func (s *Server) GenerateHtmlBodyFromTemplate(pageTitle string, pageBody string) string {
	html := s.responseTemplateFile

	html = strings.Replace(html, "{{page_title}}", pageTitle, -1)
	html = strings.Replace(html, "{{page_body}}", pageBody, -1)

	return html
}

func (s *Server) CreateResponseDetails(httpStatusCode int, reason string) (int, string, string) {
	httpStatusMessage := http.StatusText(httpStatusCode)
	htmlMessage := strings.ToLower(httpStatusMessage)

	title := httpStatusMessage

	if len(strings.TrimSpace(reason)) > 0 {
		title = fmt.Sprintf("%s.%s", httpStatusMessage, reason)
	}

	if len(strings.TrimSpace(reason)) > 0 {
		htmlMessage = s.GenerateHtmlBodyFromTemplate(title, fmt.Sprintf("HTTP Status: %d (%s, %s)", httpStatusCode, htmlMessage, reason))
	} else {
		htmlMessage = s.GenerateHtmlBodyFromTemplate(title, fmt.Sprintf("HTTP Status: %d (%s)", httpStatusCode, htmlMessage))
	}

	return httpStatusCode, httpStatusMessage, htmlMessage
}

// livenessRequestProcessor - liveness processor - is the service alive
func (s *Server) LivenessRequestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	start := time.Now().UTC()
	method := "server.livenessRequestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s entering", method)

	responseStatus := http.StatusOK
	responseMessage := "Liveness"

	s.WriteHealthCheckResponse(ctx, responseWriter, responseStatus, responseMessage)

	s.metrics.IncLivenessRequestTimer(start)
}

// readinessRequestProcessor - readiness processor
func (s *Server) ReadinessRequestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	start := time.Now().UTC()
	method := "server.readinessRequestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s entering", method)

	responseStatus := http.StatusOK
	responseMessage := "Readiness"

	if s.isShuttingDown {
		responseStatus = http.StatusServiceUnavailable
		responseMessage = "Server is shutting down."
	}

	s.WriteHealthCheckResponse(ctx, responseWriter, responseStatus, responseMessage)

	s.metrics.IncReadinessRequestTimer(start)
}

func (s *Server) WriteHealthCheckResponse(ctx context.Context, responseWriter http.ResponseWriter, httpStatusCode int, message string) {
	method := "server.writeHealthCheckHeader"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s writing response header with HTTP Status Code: %d, Message: %s", method, httpStatusCode, message)

	nodename, err := os.Hostname()
	if err == nil {
		responseWriter.Header().Add(shared.HttpHeader_Server, nodename)
	}

	responseWriter.WriteHeader(httpStatusCode)

	_, _, htmlMessage := s.CreateResponseDetails(httpStatusCode, message)

	if len(htmlMessage) > 0 {
		err = shared.WriteHTML(ctx, responseWriter, htmlMessage)
		if err != nil {
			log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseBodyContent, htmlMessage, shared.KeyErrorMessage, err.Error())).Errorf("%s error calling writeHTML", method)
		}
	}
}

// requestProcessor main server func
func (s *Server) RequestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	start := time.Now().UTC()
	method := "server.requestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s entering", method)

	if request.Method == "HEAD" { // head request responses shouldn't return a body
		s.DoHeadErrorResponse(ctx, responseWriter, request, http.StatusBadRequest, "Verb: HEAD - not supported")
		s.metrics.IncServiceRequestTimer(start)
		return
	}

	// get host name from request
	_, err := shared.GetHost(ctx, request)
	if err != nil {
		httpStatusCode, httpStatusMessage, htmlMessage := s.CreateResponseDetails(http.StatusBadRequest, err.Error())
		s.DoErrorResponse(ctx, responseWriter, request, httpStatusCode, htmlMessage, errors.New(httpStatusMessage))
		s.metrics.IncServiceRequestTimer(start)
		return
	}

	htmlMessage := fmt.Sprintf("Request from '%s' to Host: '%s', URL: '%v'", request.RemoteAddr, request.Host, request.URL)

	s.DoValidRequestResponse(ctx, responseWriter, request, htmlMessage)

	s.metrics.IncServiceRequestTimer(start)
}

// Init - setup server
func (s *Server) Init() {
	//nolint
	method := "server.Init"
	s.context = shared.CreateContext(context.Background(), s.serviceName, shared.ActionTypeService)
	log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false)).Infof("%s entering", method)

	if s.isInitialized {
		log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false)).Warnf("%s ALREADY INITIALIZED - Exiting.", method)
		return
	}

	// Init settings
	s.isInitialized = true
	s.isShuttingDown = false
	s.router = nil
	s.server = nil

	// Load Response Template HTML File
	s.responseTemplateFile = shared.LoadHTMLFile(s.context, "", DefaultResponseTemplateFile)
}

// Run - start server and listen
func (s *Server) Run() {
	//nolint
	method := "server.Run"
	log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false)).Infof("%s entering...", method)

	// setup handler
	s.router = http.NewServeMux()
	s.router.HandleFunc("/healthz/livenessZ76", s.LivenessRequestProcessor)
	s.router.HandleFunc("/healthz/readinessZ67", s.ReadinessRequestProcessor)
	s.router.HandleFunc("/", s.RequestProcessor)
	if met, ok := s.metrics.(*gometrics.GoMetrics); ok {
		s.router.Handle("/debug/gometrics", met.ExpHandler)
	}

	log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false, shared.KeyServerAddress, s.config.Service.HTTP.Server.IPv4Address)).Infof("%s server address", method)

	//tls := tls.Config{
	//}

	// NOTICE: Don't wait too long on reads or writes.
	// Holding open connections for prolonged periods is a know DDoS vector, but we have to serve
	// locations and devices with slow connections.

	// setup server
	s.server = &http.Server{
		Addr:              net.JoinHostPort(s.config.Service.HTTP.Server.IPv4Address, strconv.FormatUint(uint64(s.config.Service.HTTP.Server.Port), 10)),
		Handler:           s.router,
		ReadTimeout:       30 * time.Second, // Maximum duration for reading the entire request, including the body.
		ReadHeaderTimeout: 0,                // Amount of time allowed to read request headers. If zero, the value of ReadTimeout is used. If both are zero, there is no timeout.
		WriteTimeout:      30 * time.Second, // Maximum duration before timing out writes of the response.
		IdleTimeout:       0,                // Maximum amount of time to wait for the next request when keep-alives are enabled.  If zero, the value of ReadTimeout is used. If both are zero, there is no timeout.
		MaxHeaderBytes:    1 << 22,          // 1 << 22, 4 MB, allow for larger headers for internal users with big cookie payloads. (default: 1 << 20, aka 1 MB)
		//TLSConfig:         tls,
		//ErrorLog:          logger,
	}

	defer s.server.Close()

	err := s.server.ListenAndServe()

	s.isShuttingDown = true

	if err != nil && err != http.ErrServerClosed {
		log.WithFields(shared.GetFields(s.context, shared.EventTypeError, false, shared.KeyServerAddress, s.config.Service.HTTP.Server.IPv4Address, shared.KeyErrorMessage, err.Error())).Errorf("%s server listen error", method)
		return
	}

	log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false)).Infof("%s existing", method)
}

// Shutdown - Shutdown server
func (s *Server) Shutdown() {
	//nolint
	method := "server.Shutdown"
	log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false)).Infof("%s entering", method)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.isShuttingDown = true

	duration, err := time.ParseDuration(s.config.Service.GracefulShutdownDelaySeconds)
	if err != nil {
		duration, _ = time.ParseDuration("10s")
		log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false, shared.KeyErrorMessage, err.Error())).Errorf("%s server parsing Service.GracefulShutdownDelaySeconds, defaulting to: %s", method, duration.String())
	}

	time.Sleep(duration)

	if err = s.server.Shutdown(ctx); err != nil {
		log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false, shared.KeyErrorMessage, err.Error())).Errorf("%s server error while shutting down", method)
		panic(err)
	}
}
