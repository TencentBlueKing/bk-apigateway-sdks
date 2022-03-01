package bkapi

import (
	"encoding/json"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/internal"
	"github.com/pkg/errors"
)

// Config is the configuration of BkApi client.
type Config struct {
	apiName string

	Endpoint string
	Stage    string

	AppCode   string
	AppSecret string

	AccessToken         string
	AuthorizationParams map[string]string
	AuthorizationJWT    string
	JsonMarshaler       func(v interface{}) ([]byte, error)
}

func (c *Config) setAuthAccessTokenAuthParams(params map[string]string) bool {
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

func (c *Config) setAppAuthParams(params map[string]string) {
	if c.AppCode != "" {
		params["bk_app_code"] = c.AppCode
	}

	if c.AppSecret != "" {
		params["bk_app_secret"] = c.AppSecret
	}
}

func (c *Config) setCommonAuthParams(params map[string]string) {
	for k, v := range c.AuthorizationParams {
		params[k] = v
	}
}

func (c *Config) getAuthParams() map[string]string {
	params := make(map[string]string, 4+len(c.AuthorizationParams))

	if c.setAuthAccessTokenAuthParams(params) {
		return params
	}

	c.setAppAuthParams(params)
	c.setCommonAuthParams(params)

	return params
}

// Config method clone and return a new Config instance
func (c Config) Config(apiName string) define.ClientConfig {
	c.apiName = apiName
	return &c
}

// GetName method will return the api name.
func (c *Config) GetName() string {
	return c.apiName
}

// GetUrl method will render the endpoint with api name and stage.
func (c *Config) GetUrl() string {
	return internal.ReplacePlaceHolder(c.Endpoint, map[string]string{
		"api_name": c.apiName,
		"stage":    c.Stage,
	})
}

// GetAuthorizationHeaders method will return the authorization headers.
func (c *Config) GetAuthorizationHeaders() map[string]string {
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
