package shared_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mdonahue-godaddy/go-http-server/pkg/shared"
)

func createTestRequest(method, host, url, remoteAddr, headerKey, headerValue string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	req.Host = host
	req.RemoteAddr = remoteAddr
	if len(headerKey) > 0 {
		req.Header.Set(headerKey, headerValue)
	}
	return req
}

func createTestHttpResponseWriter() http.ResponseWriter {
	rw := httptest.NewRecorder()

	rw.Header().Add("Testing", "True")

	return rw
}

func Test_Splitter(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		Rune        rune
		Expected    bool
		Description string
	}{
		{
			Rune:        ' ',
			Expected:    true,
			Description: "' ' should return true",
		},
		{
			Rune:        ',',
			Expected:    true,
			Description: "',' should return true",
		},
		{
			Rune:        ':',
			Expected:    true,
			Description: "':' should return true",
		},
		{
			Rune:        'r',
			Expected:    false,
			Description: "'r' should return false",
		},
		{
			Rune:        't',
			Expected:    false,
			Description: "'t' should return false",
		},
	}

	for _, tc := range testCases {
		actual := shared.Splitter(tc.Rune)
		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_GetSafeString(t *testing.T) {
	assert := assert.New(t)

	empty := ""
	value := "SomeValue"

	testCases := []struct {
		Value        *string
		DefaultValue string
		Expected     string
		Description  string
	}{
		{
			Value:        nil,
			DefaultValue: "default",
			Expected:     "default",
			Description:  "nil should return default",
		},
		{
			Value:        &empty,
			DefaultValue: "default",
			Expected:     empty,
			Description:  "empty should return empty",
		},
		{
			Value:        &value,
			DefaultValue: "default",
			Expected:     value,
			Description:  "value should return value",
		},
	}

	for _, tc := range testCases {
		actual := shared.GetSafeString(tc.Value, tc.DefaultValue)
		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_GetSafeInt32AsString(t *testing.T) {
	assert := assert.New(t)

	empty := int32(0)
	emptyString := strconv.FormatInt(int64(empty), 10)
	value := int32(123)
	valueString := strconv.FormatInt(int64(value), 10)
	defaultValue := int32(999)
	defaultValueString := strconv.FormatInt(int64(defaultValue), 10)

	testCases := []struct {
		Value        *int32
		DefaultValue int32
		Expected     string
		Description  string
	}{
		{
			Value:        nil,
			DefaultValue: defaultValue,
			Expected:     defaultValueString,
			Description:  "nil should return defaultValue",
		},
		{
			Value:        &empty,
			DefaultValue: defaultValue,
			Expected:     emptyString,
			Description:  "empty should return empty",
		},
		{
			Value:        &value,
			DefaultValue: defaultValue,
			Expected:     valueString,
			Description:  "value should return value",
		},
	}

	for _, tc := range testCases {
		actual := shared.GetSafeInt32AsString(tc.Value, tc.DefaultValue)
		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_GetSafeInt64AsString(t *testing.T) {
	assert := assert.New(t)

	empty := int64(0)
	emptyString := strconv.FormatInt(empty, 10)
	value := int64(123)
	valueString := strconv.FormatInt(value, 10)
	defaultValue := int64(999)
	defaultValueString := strconv.FormatInt(defaultValue, 10)

	testCases := []struct {
		Value        *int64
		DefaultValue int64
		Expected     string
		Description  string
	}{
		{
			Value:        nil,
			DefaultValue: defaultValue,
			Expected:     defaultValueString,
			Description:  "nil should return defaultValue",
		},
		{
			Value:        &empty,
			DefaultValue: defaultValue,
			Expected:     emptyString,
			Description:  "empty should return empty",
		},
		{
			Value:        &value,
			DefaultValue: defaultValue,
			Expected:     valueString,
			Description:  "value should return value",
		},
	}

	for _, tc := range testCases {
		actual := shared.GetSafeInt64AsString(tc.Value, tc.DefaultValue)
		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_AddUniversalHeaders(t *testing.T) {
	assert := assert.New(t)

	ctx := shared.CreateContext(context.Background(), "Test_LoadHTMLFile_ActionName", "Test_LoadHTMLFile_ActionType")
	serviceName := "MyTestServiceName"
	nodename, _ := os.Hostname()

	testCases := []struct {
		Context        context.Context
		ResponseWriter http.ResponseWriter
		ServiceName    string
		Description    string
	}{
		{
			Context:        ctx,
			ResponseWriter: createTestHttpResponseWriter(),
			ServiceName:    "",
			Description:    "No service name",
		},
		{
			Context:        ctx,
			ResponseWriter: createTestHttpResponseWriter(),
			ServiceName:    serviceName,
			Description:    "Good service name",
		},
		{
			Context:        context.Background(),
			ResponseWriter: createTestHttpResponseWriter(),
			ServiceName:    serviceName,
			Description:    "Empty context",
		},
	}

	for _, tc := range testCases {
		shared.AddUniversalHeaders(tc.Context, tc.ResponseWriter, tc.ServiceName)
		assert.Equal(tc.ResponseWriter.Header().Get(shared.HttpHeader_Server), nodename, tc.Description)
		value, _ := shared.GetKeyFromContext(tc.Context, shared.KeyTransactionID)
		if value != nil {
			assert.Equal(tc.ResponseWriter.Header().Get(shared.HttpHeader_XRequestID), *value, tc.Description)
		} else {
			assert.Equal(tc.ResponseWriter.Header().Get(shared.HttpHeader_XRequestID), "", tc.Description)
		}
	}
}

func Test_WriteHTML(t *testing.T) {
	assert := assert.New(t)

	html := "<html><body>Some HTML</body></html>"

	testCases := []struct {
		ResponseWriter http.ResponseWriter
		Html           string
		ExpectedError  error
		Description    string
	}{
		{
			ResponseWriter: nil,
			Html:           html,
			ExpectedError:  errors.New("http.ResponseWriter is nil"),
			Description:    "http.ResponseWriter is nil should return error",
		},
		{
			ResponseWriter: createTestHttpResponseWriter(),
			Html:           "",
			ExpectedError:  nil,
			Description:    "http.ResponseWriter with empty html should return nil and do nothing",
		},
		{
			ResponseWriter: createTestHttpResponseWriter(),
			Html:           html,
			ExpectedError:  nil,
			Description:    "http.ResponseWriter with good html string should return nil and write the html",
		},
	}

	for _, tc := range testCases {
		ctx := shared.CreateContext(context.Background(), "Test_IsValidRequest_ActionName", "Test_IsValidRequest_ActionType")
		actual := shared.WriteHTML(ctx, tc.ResponseWriter, tc.Html)

		if tc.ExpectedError == nil {
			assert.Nil(actual, tc.Description)
		} else {
			assert.NotNil(actual, tc.Description)
			assert.Equal(tc.ExpectedError.Error(), actual.Error(), tc.Description)
		}
	}
}

func Test_GetHost(t *testing.T) {
	assert := assert.New(t)

	host := "example.com"
	url := fmt.Sprintf("http://%s/foo", host)

	testCases := []struct {
		Request        *http.Request
		ExpectedString string
		ExpectedError  error
		Description    string
	}{
		{
			Request:        nil,
			ExpectedString: "",
			ExpectedError:  errors.New("http.Request is nil or host is empty"),
			Description:    "Nil http.Request should return error",
		},
		{
			Request:        createTestRequest(http.MethodGet, "192.168.1.1", url, "", "", ""),
			ExpectedString: "192.168.1.1",
			ExpectedError:  nil,
			Description:    "http.Request with IP Address for Host should return false",
		},
		{
			Request:        createTestRequest(http.MethodGet, "IAmNot", url, "", "", ""),
			ExpectedString: "IAmNot",
			ExpectedError:  nil,
			Description:    "http.Request with invalid Host should return false FIXME",
		},
		{
			Request:        createTestRequest(http.MethodGet, host, url, "", "", ""),
			ExpectedString: host,
			ExpectedError:  nil,
			Description:    "http.Request with valid Host should return true",
		},
		{
			Request:        createTestRequest(http.MethodGet, "   :    ", url, "", "", ""),
			ExpectedString: "   ",
			ExpectedError:  errors.New("ParseUint"),
			Description:    "http.Request with space case Host should return false",
		},
	}

	for _, tc := range testCases {
		ctx := shared.CreateContext(context.Background(), "Test_IsValidRequest_ActionName", "Test_IsValidRequest_ActionType")
		actual, err := shared.GetHost(ctx, tc.Request)
		if tc.ExpectedError == nil {
			assert.Nil(err, tc.Description)
		} else {
			assert.NotNil(err, tc.Description)
		}
		assert.Equal(tc.ExpectedString, actual, tc.Description)
	}
}

func Test_IsValidGetRequest(t *testing.T) {
	assert := assert.New(t)

	host := "example.com"
	url := fmt.Sprintf("http://%s/foo", host)

	testCases := []struct {
		Request     *http.Request
		Expected    bool
		Description string
	}{
		{
			Request:     nil,
			Expected:    false,
			Description: "Nil http.Request should return false",
		},
		{
			Request:     createTestRequest(http.MethodGet, "192.168.1.1", url, "", "", ""),
			Expected:    false,
			Description: "http.Request with IP Address for Host should return false",
		},
		{
			Request:     createTestRequest(http.MethodGet, "IAmNot", url, "", "", ""),
			Expected:    true,
			Description: "http.Request with invalid Host should return false FIXME",
		},
		{
			Request:     createTestRequest(http.MethodGet, host, url, "", "", ""),
			Expected:    true,
			Description: "http.Request with valid Host should return true",
		},
		{
			Request:     createTestRequest(http.MethodPost, host, url, "", "", ""),
			Expected:    false,
			Description: "http.Request with method POST should return false",
		},
	}

	for _, tc := range testCases {
		ctx := shared.CreateContext(context.Background(), "Test_IsValidRequest_ActionName", "Test_IsValidRequest_ActionType")
		actual := shared.IsValidGetRequest(ctx, tc.Request)
		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_IsValidRequestHost(t *testing.T) {
	assert := assert.New(t)

	host := "example.com"
	url := fmt.Sprintf("http://%s/foo", host)

	testCases := []struct {
		Request     *http.Request
		Expected    bool
		Description string
	}{
		{
			Request:     nil,
			Expected:    false,
			Description: "Nil http.Request should return false",
		},
		{
			Request:     createTestRequest(http.MethodGet, "", url, "", "", ""),
			Expected:    false,
			Description: "http.Request with no Host should return false",
		},
		{
			Request:     createTestRequest(http.MethodGet, "192.168.1.1", url, "", "", ""),
			Expected:    false,
			Description: "http.Request with IP Address for Host should return false",
		},
		{
			Request:     createTestRequest(http.MethodGet, "IAmNot", url, "", "", ""),
			Expected:    true,
			Description: "http.Request with invalid Host should return false FIXME",
		},
		{
			Request:     createTestRequest(http.MethodGet, host, url, "", "", ""),
			Expected:    true,
			Description: "http.Request with valid Host should return true",
		},
		{
			Request:     createTestRequest(http.MethodGet, "   :    ", url, "", "", ""),
			Expected:    false,
			Description: "http.Request with space case Host should return false",
		},
	}

	for _, tc := range testCases {
		ctx := shared.CreateContext(context.Background(), "Test_IsValidRequestHost_ActionName", "Test_IsValidRequestHost_ActionType")
		actual := shared.IsValidRequestHost(ctx, tc.Request)
		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_IsValidHost(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		Host        string
		Expected    bool
		Description string
	}{
		{
			Host:        "",
			Expected:    false,
			Description: "Empty Host should return false",
		},
		{
			Host:        "foo",
			Expected:    true,
			Description: "Host with no TLD should return false FIXME",
		},
		{
			Host:        "foo.com",
			Expected:    true,
			Description: "Host with no TLD should return true",
		},
	}

	for _, tc := range testCases {
		ctx := shared.CreateContext(context.Background(), "Test_IsValidHost_ActionName", "Test_IsValidHost_ActionType")
		actual := shared.IsValidHost(ctx, tc.Host)
		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_LoadHTMLFile(t *testing.T) {
	assert := assert.New(t)

	DefaultHTML := "DefaultHTML"

	testCases := []struct {
		FileName     string
		DefaultHTML  string
		ExpectedHTML string
		Description  string
	}{
		{
			FileName:     "",
			DefaultHTML:  DefaultHTML,
			ExpectedHTML: DefaultHTML,
			Description:  "Empty FileName returns DefaultHTML",
		},
		{
			FileName:     "tests/test.html",
			DefaultHTML:  DefaultHTML,
			ExpectedHTML: "I am a teapot",
			Description:  "good FileName returns file contents.",
		},
		{
			FileName:     "tests/test.htmllll",
			DefaultHTML:  DefaultHTML,
			ExpectedHTML: DefaultHTML,
			Description:  "bad FileName returns DefaultHTML.",
		},
	}

	for _, tc := range testCases {
		ctx := shared.CreateContext(context.Background(), "Test_LoadHTMLFile_ActionName", "Test_LoadHTMLFile_ActionType")
		actual := shared.LoadHTMLFile(ctx, tc.FileName, tc.DefaultHTML)
		assert.Equal(tc.ExpectedHTML, actual, tc.Description)
	}
}
