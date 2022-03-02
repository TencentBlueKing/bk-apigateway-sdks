package bkapi

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gentleman.v2"
)

func newGentlemanClient(config define.ClientConfig) *gentleman.Client {
	client := gentleman.New().
		URL(config.GetUrl())

	headers := config.GetAuthorizationHeaders()
	if len(headers) > 0 {
		client.SetHeaders(headers)
	}

	return client
}

// NewBkApiClient creates a new BkApiClient.
func NewBkApiClient(apiName string, configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (define.BkApiClient, error) {
	config := configProvider.ProvideConfig(apiName)
	gentlemanClient := newGentlemanClient(config)

	client := internal.NewBkApiClient(apiName, gentlemanClient, func(name string, request *gentleman.Request) define.Operation {
		return internal.NewOperation(name, request)
	})

	if len(opts) == 0 {
		return client, nil
	}

	err := client.Apply(opts...)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to apply options to client %s", apiName)
	}

	return client, nil
}

// ClientConfig is the configuration of BkApi client.
type ClientConfig struct {
	apiName string

	Endpoint string
	Stage    string

	AppCode   string
	AppSecret string

	AccessToken         string
	AuthorizationParams map[string]string
	AuthorizationJWT    string
	JsonMarshaler       func(v interface{}) ([]byte, error)

	Getenv func(string) string
}

func (c *ClientConfig) setAuthAccessTokenAuthParams(params map[string]string) bool {
	if c.AccessToken == "" {
		return false
	}

	// if AccessToken is set, we will use it as the authorization parameters,
	// other authorization parameters will be ignored.
	params["access_token"] = c.AccessToken

	if c.AuthorizationJWT != "" {
		params["jwt"] = c.AuthorizationJWT
	}

	return true
}

func (c *ClientConfig) setAppAuthParams(params map[string]string) {
	if c.AppCode != "" {
		params["bk_app_code"] = c.AppCode
	}

	if c.AppSecret != "" {
		params["bk_app_secret"] = c.AppSecret
	}
}

func (c *ClientConfig) setCommonAuthParams(params map[string]string) {
	for k, v := range c.AuthorizationParams {
		params[k] = v
	}
}

func (c *ClientConfig) getAuthParams() map[string]string {
	params := make(map[string]string, 4+len(c.AuthorizationParams))

	if c.setAuthAccessTokenAuthParams(params) {
		return params
	}

	c.setAppAuthParams(params)
	c.setCommonAuthParams(params)

	return params
}

func (c *ClientConfig) initAppConfig() {
	if c.AppCode == "" {
		c.AppCode = c.getEnv("BK_APP_CODE", "APP_CODE")
	}

	if c.AppSecret == "" {
		c.AppSecret = c.getEnv("BK_APP_SECRET", "SECRET_KEY")
	}
}

func (c *ClientConfig) initBkApiConfig() {
	if c.Stage == "" {
		c.Stage = "prod"
	}

	endpoint := c.Endpoint
	if endpoint == "" {
		apiTmpl := c.getEnv("BK_API_URL_TMPL")
		stageTmpl := c.getEnv("BK_API_STAGE_URL_TMPL")

		if apiTmpl != "" {
			endpoint = fmt.Sprintf("%s/{stage}", strings.TrimSuffix(apiTmpl, "/"))
		} else if stageTmpl != "" {
			endpoint = stageTmpl
		}
	}

	c.Endpoint = internal.ReplacePlaceHolder(endpoint, map[string]string{
		"api_name": c.apiName,
		"stage":    c.Stage,
	})
}

func (c *ClientConfig) initConfig(apiName string) {
	c.apiName = apiName

	if c.Getenv == nil {
		c.Getenv = os.Getenv
	}

	c.initAppConfig()
	c.initBkApiConfig()
}

func (c *ClientConfig) getEnv(keys ...string) string {
	for _, k := range keys {
		v := c.Getenv(k)
		if v != "" {
			return v
		}
	}

	return ""
}

// ProvideConfig method clone and return a new Config instance
func (c ClientConfig) ProvideConfig(apiName string) define.ClientConfig {
	c.initConfig(apiName)
	return &c
}

// GetName method will return the api name.
func (c *ClientConfig) GetName() string {
	return c.apiName
}

// GetUrl method will render the endpoint with api name and stage.
func (c *ClientConfig) GetUrl() string {
	endpoint := fmt.Sprintf("%s/", strings.TrimSuffix(c.Endpoint, "/"))

	return internal.ReplacePlaceHolder(endpoint, map[string]string{
		"api_name": c.apiName,
		"stage":    c.Stage,
	})
}

// GetAuthorizationHeaders method will return the authorization headers.
func (c *ClientConfig) GetAuthorizationHeaders() map[string]string {
	params := c.getAuthParams()
	// nil means no authorization headers
	if len(params) == 0 {
		return nil
	}

	marshaler := json.Marshal
	if c.JsonMarshaler != nil {
		marshaler = c.JsonMarshaler
	}

	value, err := marshaler(params)
	if err != nil {
		// params contains basic type only, so this should never happen.
		panic(errors.WithMessagef(err, "failed to marshal bkapi authorization"))
	}

	return map[string]string{"X-Bkapi-Authorization": string(value)}
}
