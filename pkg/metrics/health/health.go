package health

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const (
	HttpHeader_ContentType string = "Content-Type"
	ContentType_TextHtml   string = "text/html; charset=utf-8"

	DefaultResponseTemplate string = "<!DOCTYPE HTML><html lang='en-us'><head><title>** {{page_title}} **</title></head><body>{{page_body}}</body></html>"
	Status_Good             string = "GOOD!"
)

type Status struct {
	IsGood bool
	Reason string
}

type HealthEndpoints struct {
	ServeMux *http.ServeMux
	Base     string
	Status   Status
}

func NewHealthEndpoints(mux *http.ServeMux, base string) *HealthEndpoints {
	endpoint := HealthEndpoints{
		ServeMux: mux,
		Base:     base,
		Status: Status{
			IsGood: false,
			Reason: "Status now set yet.",
		},
	}

	return &endpoint
}

func (s *HealthEndpoints) SetStatus(isGood bool, reason string) {
	s.Status.IsGood = isGood
	s.Status.Reason = reason
}

// EnableEndpoints - add liveness and readiness handlers
func (s *HealthEndpoints) EnableEndpoints() {
	s.ServeMux.HandleFunc(fmt.Sprintf("%s/liveness", s.Base), s.LivenessHandler)
	s.ServeMux.HandleFunc(fmt.Sprintf("%s/readiness", s.Base), s.ReadinessHandler)
}

// LivenessHandler - liveness handler, service is alive
func (s *HealthEndpoints) LivenessHandler(responseWriter http.ResponseWriter, request *http.Request) {
	htmlMessage := strings.Replace(DefaultResponseTemplate, "{{page_title}}", "Liveness Status", 1)

	if s.Status.IsGood {
		responseWriter.WriteHeader(http.StatusOK)

		htmlMessage = strings.Replace(htmlMessage, "{{page_body}}", Status_Good, 1)
	} else {
		responseWriter.WriteHeader(http.StatusOK)

		htmlMessage = strings.Replace(htmlMessage, "{{page_body}}", Status_Good, 1)
	}

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

// ReadinessHandler - readiness handler, service is ready
func (s *HealthEndpoints) ReadinessHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
