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

	"github.com/mdonahue-godaddy/go-http-server/pkg/config"
	"github.com/mdonahue-godaddy/go-http-server/pkg/shared"
)

const (
	defaultResponseTemplateFile string = "<!DOCTYPE HTML><html lang='en-us'><head><title>** {{page_title}} **</title></head><body>{{page_body}}</body></html>"
)

// ForwardingHTTPResponse struct
type ForwardingHTTPResponse struct {
	RequestedHost string
	Error         error
}

// Server struct
type Server struct {
	serviceName          string
	isInitialized        bool
	isShuttingDown       bool
	config               *config.Settings
	context              context.Context
	router               *http.ServeMux
	server               *http.Server
	responseTemplateFile string
}

// NewServer - create new instance of server
func NewServer(serviceName string, cfg *config.Settings) *Server {
	server := new(Server)

	server.serviceName = serviceName
	server.config = cfg
	server.responseTemplateFile = defaultResponseTemplateFile

	return server
}

// doHeadErrorResponse - WARNING HEAD requests should not return a body so normal error response can't be used.
func (s *Server) doHeadErrorResponse(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, httpStatusCode int, message string) {
	method := "server.doHeadRequestErrorResponse"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseStatusCode, httpStatusCode)).Debugf("%s entering", method)

	responseWriter.Header().Add("X-Error", fmt.Sprintf("HEAD not supported, %s", message))
	responseWriter.Header().Add("Content-Length", "0")
	responseWriter.Header().Add("Content-Type", shared.ContentType_TextHtml)

	s.writeHeader(ctx, responseWriter, httpStatusCode)
}

// doErrorResponse - Error response processor
func (s *Server) doErrorResponse(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, httpStatusCode int, htmlMessage string, err error) {
	method := "server.doErrorResponse"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseStatusCode, httpStatusCode, shared.KeyHTTPResponseBodyContent, htmlMessage)).Debugf("%s entering", method)

	shared.AddUniversalHeaders(ctx, responseWriter, s.serviceName)

	msg := ""
	if err != nil {
		msg = err.Error()
	}

	log.WithFields(shared.GetFields(ctx, shared.EventTypeError, false, shared.KeyHTTPResponseStatusCode, httpStatusCode, shared.KeyHTTPResponseBodyContent, htmlMessage)).Errorf("%s returning error response. %s", method, msg)

	if request.Method == "HEAD" { // head request responses shouldn't return a body
		s.doHeadErrorResponse(ctx, responseWriter, request, httpStatusCode, msg)
		return
	}

	s.writeHeader(ctx, responseWriter, httpStatusCode)

	if len(htmlMessage) > 0 {
		err := shared.WriteHTML(ctx, responseWriter, htmlMessage)
		if err != nil {
			log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseStatusCode, httpStatusCode, shared.KeyHTTPResponseBodyContent, htmlMessage, shared.KeyErrorMessage, err.Error())).Errorf("%s error calling writeHTML", method)
		}
	}
}

func (s *Server) doValidRequestResponse(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, responseHTML string) {
	method := "server.doValidRequestResponse"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s entering", method)

	shared.AddUniversalHeaders(ctx, responseWriter, s.serviceName)

	// Http Status Code 200
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s Masking", method)

	s.writeHeader(ctx, responseWriter, http.StatusOK)

	err := shared.WriteHTML(ctx, responseWriter, responseHTML)
	if err != nil {
		log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyErrorMessage, err.Error())).Errorf("%s writeHTML error", method)
	}
}

func (s *Server) writeHeader(ctx context.Context, responseWriter http.ResponseWriter, httpStatusCode int) {
	method := "server.writeHeader"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s writing response header with HTTP Status Code: %d", method, httpStatusCode)

	responseWriter.WriteHeader(httpStatusCode)
}

func (s *Server) generateHtmlBodyFromTemplate(pageTitle string, pageBody string) string {
	html := s.responseTemplateFile

	html = strings.Replace(html, "{{page_title}}", pageTitle, -1)
	html = strings.Replace(html, "{{page_body}}", pageBody, -1)

	return html
}

func (s *Server) createResponseDetails(httpStatusCode int, reason string) (int, string, string) {
	httpStatusMessage := http.StatusText(httpStatusCode)
	htmlMessage := strings.ToLower(httpStatusMessage)

	title := httpStatusMessage

	if len(strings.TrimSpace(reason)) > 0 {
		title = fmt.Sprintf("%s.%s", httpStatusMessage, reason)
	}

	if len(strings.TrimSpace(reason)) > 0 {
		htmlMessage = s.generateHtmlBodyFromTemplate(title, fmt.Sprintf("HTTP Status: %d (%s, %s)", httpStatusCode, htmlMessage, reason))
	} else {
		htmlMessage = s.generateHtmlBodyFromTemplate(title, fmt.Sprintf("HTTP Status: %d (%s)", httpStatusCode, htmlMessage))
	}

	return httpStatusCode, httpStatusMessage, htmlMessage
}

// livenessRequestProcessor - liveness processor - is the service alive
func (s *Server) livenessRequestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	method := "server.livenessRequestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s entering", method)

	s.writeHealthCheckResponse(ctx, responseWriter, http.StatusOK, "Liveness")
}

// readinessRequestProcessor - readiness processor
func (s *Server) readinessRequestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	method := "server.readinessRequestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s entering", method)

	if s.isShuttingDown {
		s.writeHealthCheckResponse(ctx, responseWriter, http.StatusServiceUnavailable, "Server is shutting down.")
		return
	}

	s.writeHealthCheckResponse(ctx, responseWriter, http.StatusOK, "Readiness")
}

func (s *Server) writeHealthCheckResponse(ctx context.Context, responseWriter http.ResponseWriter, httpStatusCode int, message string) {
	method := "server.writeHealthCheckHeader"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s writing response header with HTTP Status Code: %d, Message: %s", method, httpStatusCode, message)

	nodename, err := os.Hostname()
	if err == nil {
		responseWriter.Header().Add(shared.HttpHeader_Server, nodename)
	}

	responseWriter.WriteHeader(httpStatusCode)

	_, _, htmlMessage := s.createResponseDetails(httpStatusCode, message)

	if len(htmlMessage) > 0 {
		err = shared.WriteHTML(ctx, responseWriter, htmlMessage)
		if err != nil {
			log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false, shared.KeyHTTPResponseBodyContent, htmlMessage, shared.KeyErrorMessage, err.Error())).Errorf("%s error calling writeHTML", method)
		}
	}
}

// requestProcessor main server func
func (s *Server) requestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	method := "server.requestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s entering", method)

	// get host name from request
	requestedHost, err := shared.GetHost(ctx, request)
	if err != nil {
		httpStatusCode, httpStatusMessage, htmlMessage := s.createResponseDetails(http.StatusBadRequest, err.Error())
		s.doErrorResponse(ctx, responseWriter, request, httpStatusCode, htmlMessage, errors.New(httpStatusMessage))
	}

	// Init variables
	httpStatusCode := http.StatusOK
	httpStatusMessage := http.StatusText(httpStatusCode)
	htmlMessage := fmt.Sprintf("From Host: '%s'", requestedHost)

	if shared.IsValidGetRequest(ctx, request) {
		// success
		s.doValidRequestResponse(ctx, responseWriter, request, htmlMessage)

		return
	}

	// Not Valid GET Request
	httpStatusCode, httpStatusMessage, htmlMessage = s.createResponseDetails(http.StatusBadRequest, "Request Not Valid")
	s.doErrorResponse(ctx, responseWriter, request, httpStatusCode, htmlMessage, errors.New(httpStatusMessage))
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
	s.responseTemplateFile = shared.LoadHTMLFile(s.context, "", defaultResponseTemplateFile)
}

// Run - start server and listen
func (s *Server) Run() {
	//nolint
	method := "server.Run"
	log.WithFields(shared.GetFields(s.context, shared.EventTypeInfo, false)).Infof("%s entering...", method)

	// setup handler
	s.router = http.NewServeMux()
	s.router.HandleFunc("/livenessZ76", s.livenessRequestProcessor)
	s.router.HandleFunc("/readinessZ67", s.readinessRequestProcessor)
	s.router.HandleFunc("/", s.requestProcessor)

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