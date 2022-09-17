package metrics

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/pprof"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mdonahue-godaddy/go-http-server/pkg/shared"
)

const (
	ContentType_TextHtml    string = "text/html; charset=utf-8"
	HttpHeader_ContentType  string = "Content-Type"
	DefaultResponseTemplate string = "<!DOCTYPE HTML><html lang='en-us'><head><title>** {{page_title}} **</title></head><body>{{page_body}}</body></html>"
	Status_Good             string = "GOOD!"
)

var (
	DefaultDialAddress = ":8082"
)

// StartServer is the entry point into initializing a pprof server instance for the Features API
func StartServer(ctx context.Context, dialAddress string, enableHealthCheck bool, enablePProf bool) {
	method := "metrics.StartServer"

	var err error

	if enableHealthCheck || enablePProf {
		if dialAddress == "" {
			dialAddress = DefaultDialAddress
		}

		log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s starting http server on Dial Address: %s", method, dialAddress)

		mux := http.NewServeMux()

		if enableHealthCheck {
			err = StartHealthCheck(ctx, mux, "/healthz")
			if err != nil {
				log.WithFields(shared.GetFields(ctx, shared.EventTypeError, false, shared.KeyErrorMessage, err.Error())).Errorf("%s error calling StartHealthCheck()", method)
			}
		}

		if enablePProf {
			err = StartPProfAPI(ctx, mux, "/debugz")
			if err != nil {
				log.WithFields(shared.GetFields(ctx, shared.EventTypeError, false, shared.KeyErrorMessage, err.Error())).Errorf("%s error calling StartPProfAPI()", method)
			}
		}

		if err := http.ListenAndServe(dialAddress, mux); err != nil {
			log.WithFields(shared.GetFields(ctx, shared.EventTypeError, false, shared.KeyErrorMessage, err.Error())).Errorf("%s error calling http.ListenAndServe() for http server on Dial Address: %s", method, dialAddress)
		}
	}
}

// StartHealthCheck starts pprof endpoint
func StartHealthCheck(ctx context.Context, mux *http.ServeMux, base string) error {
	method := "service.StartHealthCheck"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s starting endpoint", method)

	mux.HandleFunc(fmt.Sprintf("%s/liveness", base), livenessRequestProcessor)
	mux.HandleFunc(fmt.Sprintf("%s/readiness", base), readinessRequestProcessor)

	return nil
}

// livenessRequestProcessor - liveness processor - is the service alive
func livenessRequestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	method := "service.livenessRequestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s starting endpoint", method)

	responseWriter.WriteHeader(http.StatusOK)

	htmlMessage := strings.Replace(DefaultResponseTemplate, "{{page_title}}", "Liveness Status", 1)
	htmlMessage = strings.Replace(htmlMessage, "{{page_body}}", Status_Good, 1)

	if len(htmlMessage) > 0 {
		responseWriter.Header().Set(HttpHeader_ContentType, ContentType_TextHtml)
		msg := template.New("body")
		var err error
		msg, err = msg.Parse(htmlMessage)
		if err != nil {
			return
		}

		_ = msg.Execute(responseWriter, nil)
	}
}

// readinessRequestProcessor - readiness processor
func readinessRequestProcessor(responseWriter http.ResponseWriter, request *http.Request) {
	method := "service.readinessRequestProcessor"
	ctx := shared.CreateRequestContext(request, method)
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Debugf("%s starting endpoint", method)

	responseWriter.WriteHeader(http.StatusOK)

	htmlMessage := strings.Replace(DefaultResponseTemplate, "{{page_title}}", "Readiness Status", 1)
	htmlMessage = strings.Replace(htmlMessage, "{{page_body}}", Status_Good, 1)

	if len(htmlMessage) > 0 {
		responseWriter.Header().Set(HttpHeader_ContentType, ContentType_TextHtml)
		msg := template.New("body")
		var err error
		msg, err = msg.Parse(htmlMessage)
		if err != nil {
			return
		}

		_ = msg.Execute(responseWriter, nil)
	}
}

// StartPProfAPI starts pprof endpoint
func StartPProfAPI(ctx context.Context, mux *http.ServeMux, base string) error {
	method := "service.StartPProfAPI"
	log.WithFields(shared.GetFields(ctx, shared.EventTypeInfo, false)).Infof("%s starting endpoint", method)

	mux.HandleFunc(fmt.Sprintf("%s/pprof/", base), pprof.Index)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/cmdline", base), pprof.Cmdline)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/profile", base), pprof.Profile)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/symbol", base), pprof.Symbol)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/trace", base), pprof.Trace)

	mux.Handle(fmt.Sprintf("%s/pprof/goroutine", base), pprof.Handler("goroutine"))
	mux.Handle(fmt.Sprintf("%s/pprof/heap", base), pprof.Handler("heap"))
	mux.Handle(fmt.Sprintf("%s/pprof/threadcreate", base), pprof.Handler("threadcreate"))
	mux.Handle(fmt.Sprintf("%s/pprof/block", base), pprof.Handler("block"))
	mux.Handle(fmt.Sprintf("%s/vars", base), http.DefaultServeMux)

	return nil
}
