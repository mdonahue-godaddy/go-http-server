package shared

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// ContextKey comprises keys used to access common information from a request context.
type ContextKey string

// ContextValuesKey - ...
type ContextValuesKey string

var (
	config map[string]interface{} = make(map[string]interface{})
)

const (
	// HTTPRequestActionType - HTTPRequest
	HTTPRequestActionType string = "HTTPRequest"

	// ValuesKey is ...
	ValuesKey ContextValuesKey = "Values"

	// KeyActionName is ...
	KeyActionName string = "action.name"
	// KeyActionType is ...
	KeyActionType string = "action.type"
	// KeyActionStart is ...
	KeyActionStart string = "action.start"

	// KeyEventKind is required, allowed values: alert, event, metric, state, signal -- NOTICE: pipeline_error is also allowed in the field, but it's only allow to be used by elasticsearch
	KeyEventKind string = "event.kind"
	// KeyEventCategory is required, allowed values: authentication, configuration, database, driver, file, host, iam, intrusion_detection, malware, network, package, process, web
	KeyEventCategory string = "event.category"
	// KeyEventType is required, allowed values: access, admin, allowed, change, connection, creation, deletion, denied, end, error, group, info, installation, protocol, start, user
	KeyEventType string = "event.type"
	// KeyEventOutcome is optional, allowed values: failure, success, unknown
	KeyEventOutcome string = "event.outcome"
	// KeyCloudAccountID is required for AWS
	KeyCloudAccountID string = "cloud.account.id"
	// KeyCloudAccountName is
	KeyCloudAccountName string = "cloud.account.name"
	// KeyCloudProvider is ...
	KeyCloudProvider string = "cloud.provider"
	// KeyCloudRegion is ...
	KeyCloudRegion string = "cloud.region"
	// KeyCloudInstanceID is ...
	KeyCloudInstanceID string = "cloud.instance.id"
	// KeyCloudInstanceName is required for Lambda
	KeyCloudInstanceName string = "cloud.instance.name"
	// KeyTimestamp is ...
	KeyTimestamp string = "@timestamp"
	// KeyTag is a required array with 'security', 'application' or both
	KeyTag string = "tag"
	// KeyServiceName (If AWS, prefix with AWS account name ..i.e. domains-sellerapi)
	KeyServiceName string = "service.name"
	// KeyTransactionID is used for tieing together multiple related events
	KeyTransactionID string = "transaction.id"
	// KeyTraceID is required if available
	KeyTraceID string = "trace.id"
	// KeySourceAddress is required if available
	KeySourceAddress string = "source.address"
	// KeyServerAddress is
	KeyServerAddress string = "server.address"

	// KeyRequestProto is ...
	KeyRequestProto string = "request.proto"
	// KeyRequestMethod is ...
	KeyRequestMethod string = "request.method"
	// KeyRequestHost is ...
	KeyRequestHost string = "request.host"
	// KeyRequestURL is ...
	KeyRequestURL string = "request.url"
	// KeyRequestURI is ...
	KeyRequestURI string = "request.uri"
	// KeyRequestRemoteAddr is ...
	KeyRequestRemoteAddr string = "request.remoteaddr" // should use source address
	// KeyDBRetries is max number of database connect retries
	KeyDBRetries string = "db.retries"
	// KeyDBCurrentTryCount is Current Try/Retry Count
	KeyDBCurrentTryCount string = "db.current_try_count"
	// KeyHTTPResponseBodyContent is HTTP Response Body Content
	KeyHTTPResponseBodyContent string = "http.response.body.content"
	// KeyHTTPResponseStatusCode is HTTP Response Status Code
	KeyHTTPResponseStatusCode string = "http.response.status_code"
	// KeyMapEntryType is Map Entry Type
	KeyMapEntryType string = "map_entry_type"
	// KeyProcessLocalConfig is Process Local Configuration File
	KeyProcessLocalConfig string = "process.local.config"
	// KeyProcessConfigEnv is Process Configured Environment
	KeyProcessConfigEnv string = "process.config.env"
	// KeyProcessConfigJSON is Process Configuration JSON
	KeyProcessConfigJSON string = "process.config.json"
	// KeyMaskHeaderURL is Mask Header URL
	KeyMaskHeaderURL string = "mask_header_url"
	//KeyAppConfig is Key to AppConfig
	KeyAppConfig string = "AppConfig"
	//KeyAppConfigKey is Key to AppConfigKey
	KeyAppConfigKey string = "AppConfigKey"

	// KeyErrorMessage is error.message
	KeyErrorMessage string = "error.message"
	// KeyErrorCode is error.code
	KeyErrorCode string = "error.code"

	// KeyArgs is args
	KeyArgs string = "args"
	// KeyCode is code
	KeyCode string = "code"
	// KeyDomainName is domain name
	KeyDomainName string = "domain_name"
	// KeyQuery is query
	KeyQuery string = "query"
	// KeyShopperID i shopper_id
	KeyShopperID string = "shopper_id"

	// Valid EventType values: access, admin, allowed, change, connection, creation, deletion, denied, end, error, group, info, installation, protocol, start, user

	// EventTypeError is error
	EventTypeError = "error"
	// EventTypeInfo is info
	EventTypeInfo = "info"

	// ActionTypeService is service
	ActionTypeService = "service"
)

// Init - set config values
func Init(serviceName string) {
	config[KeyCloudInstanceName], _ = os.Hostname()
	config[KeyServiceName] = serviceName
}

func getConfigMap() map[string]interface{} {
	return config
}

// Values - struct to store values on context
type Values struct {
	m map[string]interface{}
}

// Get - Get value from Values struct
func (v Values) Get(key string) interface{} {
	return v.m[key]
}

// CreateContext - ....
func CreateContext(ctx context.Context, actionName string, actionType string) context.Context {
	values := Values{map[string]interface{}{
		KeyActionName:    actionName,
		KeyActionType:    actionType,
		KeyActionStart:   time.Now().UTC().Format(TimestampFormat),
		KeyTransactionID: uuid.New().String(),
	}}

	return context.WithValue(ctx, ValuesKey, values)
}

// CreateRequestContext - ....
func CreateRequestContext(request *http.Request, actionName string) context.Context {
	values := Values{map[string]interface{}{
		KeyActionName:        actionName,
		KeyActionType:        HTTPRequestActionType,
		KeyActionStart:       time.Now().UTC().Format(TimestampFormat),
		KeyTransactionID:     uuid.New().String(),
		KeyRequestProto:      request.Proto,
		KeyRequestMethod:     request.Method,
		KeyRequestHost:       request.Host,
		KeyRequestURL:        request.URL.String(),
		KeyRequestURI:        request.RequestURI,
		KeyRequestRemoteAddr: request.RemoteAddr,
	}}

	return context.WithValue(request.Context(), ValuesKey, values)
}

// GetFields - get
func GetFields(ctx context.Context, eventType string, addSecurityTag bool, fields ...interface{}) map[string]interface{} {
	results := make(map[string]interface{})

	configFields := getConfigMap()
	for k, v := range configFields {
		results[k] = v
	}

	contextFields := getContextMap(ctx)
	for k, v := range contextFields {
		results[k] = v
	}

	if len(fields) > 0 {
		fm := pairFields(fields...)
		for k, v := range fm {
			results[k] = v
		}
	}

	results[KeyEventType] = eventType

	if addSecurityTag {
		results[KeyTag] = []string{"application", "security"}
	} else {
		results[KeyTag] = []string{"application"}
	}

	return results
}

func getContextMap(ctx context.Context) map[string]interface{} {
	results := make(map[string]interface{})
	values := ctx.Value(ValuesKey)

	if values != nil {
		actionType := values.(Values).Get(KeyActionType)

		results[KeyActionName] = values.(Values).Get(KeyActionName)
		results[KeyActionType] = actionType
		results[KeyActionStart] = values.(Values).Get(KeyActionStart)
		results[KeyTransactionID] = values.(Values).Get(KeyTransactionID)

		if actionType == HTTPRequestActionType {
			results[KeyRequestProto] = values.(Values).Get(KeyRequestProto)
			results[KeyRequestMethod] = values.(Values).Get(KeyRequestMethod)
			results[KeyRequestHost] = values.(Values).Get(KeyRequestHost)
			results[KeyRequestURL] = values.(Values).Get(KeyRequestURL)
			results[KeyRequestURI] = values.(Values).Get(KeyRequestURI)
			results[KeyRequestRemoteAddr] = values.(Values).Get(KeyRequestRemoteAddr)
		}
	}

	return results
}

// GetKeyFromContext get value from Values map stored in context
func GetKeyFromContext(ctx context.Context, key string) (*string, error) {
	values := ctx.Value(ValuesKey)

	if values == nil {
		return nil, errors.New("missing ValuesKey")
	}

	ivalue := values.(Values).Get(key)

	if ivalue == nil {
		return nil, errors.New("key not found")
	}

	switch value := ivalue.(type) {
	case string:
		return &value, nil
	default:
		return nil, errors.New("value not string")
	}
}

func pairFields(fields ...interface{}) map[string]interface{} {
	results := make(map[string]interface{})

	end := len(fields)

	for idx := 0; idx < end; {
		key := fields[idx].(string)

		var value interface{}

		if (idx + 1) < end {
			value = fields[idx+1]
		}

		switch valueType := value.(type) {
		case []interface{}:
			value = pairFields(valueType...)
		}

		results[key] = value

		idx += 2
	}

	return results
}
