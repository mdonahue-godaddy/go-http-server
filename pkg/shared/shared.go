package shared

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

const (
	ContentType_TextHtml string = "text/html; charset=utf-8"

	HttpHeader_ContentType = "Content-Type"
	HttpHeader_Server      = "Server"
	HttpHeader_XRequestID  = "X-Request-ID"
)

func Splitter(r rune) bool {
	return r == ' ' || r == ',' || r == ':'
}

func GetSafeString(value *string, defaultValue string) string {
	if value == nil {
		return defaultValue
	}

	return *value
}

func GetSafeInt32AsString(value *int32, defaultValue int32) string {
	if value == nil {
		return strconv.FormatInt(int64(defaultValue), 10)
	}

	return strconv.FormatInt(int64(*value), 10)
}

func GetSafeInt64AsString(value *int64, defaultValue int64) string {
	if value == nil {
		return strconv.FormatInt(defaultValue, 10)
	}

	return strconv.FormatInt(*value, 10)
}

func AddUniversalHeaders(ctx context.Context, responseWriter http.ResponseWriter, serviceName string) {
	method := "shared.AddUniversalHeaders"
	log.WithFields(GetFields(ctx, EventTypeInfo, false)).Debugf("%s entering", method)

	nodename, err := os.Hostname()
	if err == nil {
		responseWriter.Header().Add(HttpHeader_Server, nodename)
	} else {
		responseWriter.Header().Add(HttpHeader_Server, serviceName)
		log.WithFields(GetFields(ctx, EventTypeError, false, KeyErrorMessage, err.Error())).Warnf("%s os.Hostname() error", method)
	}

	requestID, err := GetKeyFromContext(ctx, KeyTransactionID)
	if err == nil {
		responseWriter.Header().Add(HttpHeader_XRequestID, *requestID)
	} else {
		log.WithFields(GetFields(ctx, EventTypeInfo, false, KeyErrorMessage, err.Error())).Warnf("%s requestID error", method)
	}
}

func WriteHTML(ctx context.Context, responseWriter http.ResponseWriter, html string) error {
	method := "shared.WriteHTML"
	log.WithFields(GetFields(ctx, EventTypeInfo, false, KeyHTTPResponseBodyContent, html)).Debugf("%s entering", method)

	var err error

	if responseWriter == nil {
		return errors.New("http.ResponseWriter is nil")
	}

	if len(html) > 0 {
		responseWriter.Header().Set(HttpHeader_ContentType, ContentType_TextHtml)

		msg := template.New("body")
		msg, err = msg.Parse(html)
		if err != nil {
			return err
		}

		err = msg.Execute(responseWriter, nil)
	}

	return err
}

// SplitHost - Split HTTPRequest.Host or HTTPRequest.RemoteHost into host and port
func SplitHost(host string) (string, uint64, error) {
	err := error(nil)
	hostName := host
	uport := uint64(0)

	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		hostName = parts[0]

		uport, err = strconv.ParseUint(parts[1], 10, 64)
	}

	return hostName, uport, err
}

func GetHost(ctx context.Context, request *http.Request) (string, error) {
	method := "shared.GetHost"
	log.WithFields(GetFields(ctx, EventTypeInfo, false)).Debugf("%s entering", method)

	if (request == nil) || (len(request.Host) <= 0) {
		err := errors.New("http.Request is nil or host is empty")
		log.WithFields(GetFields(ctx, EventTypeInfo, false, KeyErrorMessage, err.Error())).Errorf("%s request error", method)
		return "", err
	}

	requestedHost, _, err := SplitHost(request.Host)
	if err != nil {
		log.WithFields(GetFields(ctx, EventTypeInfo, false, KeyErrorMessage, err.Error())).Errorf("%s SplitHost error", method)
	}

	return requestedHost, err
}

func IsValidGetRequest(ctx context.Context, request *http.Request) bool {
	method := "shared.IsValidRequest"
	log.WithFields(GetFields(ctx, EventTypeInfo, false)).Debugf("%s entering", method)

	if request == nil {
		log.WithFields(GetFields(ctx, EventTypeError, false)).Errorf("%s http.Request is nil.", method)
		return false
	}

	if request.Method == http.MethodGet {
		if IsValidRequestHost(ctx, request) {
			return true
		}

		log.WithFields(GetFields(ctx, EventTypeError, false, KeyRequestHost, request.Host)).Warnf("%s invalid request.Host", method)
	} else {
		log.WithFields(GetFields(ctx, EventTypeError, false, KeyRequestMethod, request.Method)).Warnf("%s invalid request.Method", method)
	}

	return false
}

// IsIPAddress will check if an ipaddress is valid v4 or v6 address
func IsIPAddress(ipAddress string) bool {
	parsed := net.ParseIP(ipAddress)

	if parsed == nil {
		return false
	}

	if parsed.To4() != nil {
		return true
	}

	if parsed.To16() != nil {
		return true
	}

	return false
}

func IsValidRequestHost(ctx context.Context, request *http.Request) bool {
	method := "shared.IsValidRequestHost"
	log.WithFields(GetFields(ctx, EventTypeInfo, false)).Debugf("%s entering", method)

	if request == nil {
		log.WithFields(GetFields(ctx, EventTypeError, false, KeyErrorMessage, "http.Request is nil")).Warnf("%s http.Request is nil", method)
		return false
	} else if len(request.Host) <= 0 {
		log.WithFields(GetFields(ctx, EventTypeError, false, KeyErrorMessage, "http.Request.Host is empty")).Warnf("%s http.Request.Host is empty", method)
		return false
	}

	requestedHost, _, err := SplitHost(request.Host)
	if err != nil {
		log.WithFields(GetFields(ctx, EventTypeError, false, KeyRequestHost, request.Host, KeyErrorMessage, err.Error())).Warnf("%s SplitHost(%s) returned error", method, request.Host)
		return false
	}

	if IsIPAddress(requestedHost) {
		log.WithFields(GetFields(ctx, EventTypeError, false, KeyRequestHost, requestedHost)).Warnf("%s IsIPAddress(%s) is true, IP Address in host not allowed", method, requestedHost)
		return false
	}

	return IsValidHost(ctx, requestedHost)
}

func IsValidHost(ctx context.Context, host string) bool {
	method := "shared.IsValidHost"
	log.WithFields(GetFields(ctx, EventTypeInfo, false)).Debugf("%s entering", method)

	result := govalidator.IsDNSName(host)
	return result
}

func LoadHTMLFile(ctx context.Context, fileName string, defaultHTML string) string {
	method := "shared.LoadHTMLFile"
	log.WithFields(GetFields(ctx, EventTypeInfo, false)).Debugf("%s entering", method)

	if len(fileName) > 0 {
		// Read HTTML File
		buffer, err := ioutil.ReadFile(fileName)

		if err != nil {
			log.WithFields(GetFields(ctx, EventTypeError, false, "fileName", fileName, KeyErrorMessage, err.Error())).Warnf("%s error loading HTML from file, using hard coded default", method)
			return defaultHTML
		}

		return string(buffer)
	}

	return defaultHTML
}

// Struct2JSONString converts struct to JSON string, typically for logging
func Struct2JSONString(value interface{}) (*string, error) {
	bytes, err := json.Marshal(value)

	if err != nil {
		return nil, err
	}

	jsonString := string(bytes)

	return &jsonString, nil
}
