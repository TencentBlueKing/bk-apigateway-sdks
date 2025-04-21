/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package bkapi

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TencentBlueKing/gopkg/logging"
	gentleman "gopkg.in/h2non/gentleman.v2"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/internal"
)

// the common options for all bkapi clients
var globalBkapiClientOptions []define.BkApiClientOption

// RegisterGlobalBkapiClientOption use to register a global bkapi client option.
// Warning: this function is not safe for concurrent access.
func RegisterGlobalBkapiClientOption(opt define.BkApiClientOption) {
	globalBkapiClientOptions = append(globalBkapiClientOptions, opt)
}

// NewBkApiClient creates a new BkApiClient.
func NewBkApiClient(
	apiName string,
	configProvider define.ClientConfigProvider,
	options ...define.BkApiClientOption,
) (define.BkApiClient, error) {
	config := configProvider.ProvideConfig(apiName)
	client, err := internal.NewBkApiClient(
		apiName,
		gentleman.New(),
		func(name string, client define.BkApiClient, request *gentleman.Request) define.Operation {
			return internal.NewOperation(name, client, request)
		},
		config,
	)
	if err != nil {
		return nil, err
	}

	for phase, opts := range map[string][]define.BkApiClientOption{
		"global": globalBkapiClientOptions,
		"config": config.GetClientOptions(),
		"client": options,
	} {
		if len(opts) == 0 {
			continue
		}

		err := client.Apply(opts...)
		if err != nil {
			return nil, define.ErrorWrapf(err, "failed to apply options to client %s, phase %s", apiName, phase)
		}
	}

	return client, nil
}

// ClientConfig is the configuration of BkApi client.
type ClientConfig struct {
	apiName string

	// Endpoint is the url of the BkApi server.
	// Default: "{BkApiUrlTmpl}/{Stage}"
	Endpoint string
	// BkApiUrlTmpl is the template for endpoint,
	// Default: os.Getenv("BK_API_URL_TMPL")
	BkApiUrlTmpl string
	// Stage is the api stage name, defaults to "prod".
	// Default: "prod"
	Stage string

	// AppCode is the blueking app code.
	AppCode string
	// AppSecret is the secret key of the blueking app.
	AppSecret string

	// AccessToken is the access token of the user and app, optional.
	AccessToken string
	// AuthorizationParams is the authorization params of the user and app, optional.
	AuthorizationParams map[string]string
	// AuthorizationJWT is the bkapi jwt, optional.
	AuthorizationJWT string
	// JsonMarshal is the json marshal function, defaults to json.Marshal.
	JsonMarshaler func(v interface{}) ([]byte, error)

	// Getenv is the function to get env, defaults to os.Getenv.
	Getenv func(string) string

	// Logger is used to log the request and response.
	Logger logging.Logger

	// ClientConfig will apply to the client.
	ClientOptions []define.BkApiClientOption
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
		c.AppCode = c.getEnv("BK_APP_CODE", "BKPAAS_APP_ID", "APP_CODE")
	}

	if c.AppSecret == "" {
		c.AppSecret = c.getEnv("BK_APP_SECRET", "BKPAAS_APP_SECRET", "SECRET_KEY")
	}
}

func (c *ClientConfig) initBkApiConfig() {
	if c.Endpoint != "" {
		return
	}

	if c.Stage == "" {
		c.Stage = "prod"
	}

	if c.BkApiUrlTmpl == "" {
		c.BkApiUrlTmpl = c.getEnv("BK_API_URL_TMPL")
	}

	var endpoint string
	if c.BkApiUrlTmpl != "" {
		endpoint = fmt.Sprintf("%s/{stage}", strings.TrimSuffix(c.BkApiUrlTmpl, "/"))
	} else {
		endpoint = c.getEnv("BK_API_STAGE_URL_TMPL")
	}

	c.Endpoint = internal.ReplacePlaceHolder(endpoint, map[string]string{
		"api_name": c.apiName,
		"stage":    c.Stage,
	})
}

func (c *ClientConfig) initLogger() {
	if c.Logger != nil {
		return
	}

	c.Logger = logging.GetLogger(loggerName)
}

func (c *ClientConfig) setApiName(apiName string) {
	c.apiName = apiName
}

func (c *ClientConfig) initConfig(apiName string) {
	c.setApiName(apiName)

	if c.Getenv == nil {
		c.Getenv = os.Getenv
	}

	c.initAppConfig()
	c.initBkApiConfig()
	c.initLogger()
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

// ProvideConfig method clone and return a new Config instance.
// This method should fill the default values which are not set.
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
		panic(define.ErrorWrapf(err, "failed to marshal bkapi authorization"))
	}

	return map[string]string{"X-Bkapi-Authorization": string(value)}
}

// GetLogger method will return the logger.
func (c *ClientConfig) GetLogger() logging.Logger {
	return c.Logger
}

// GetClientOptions method will return the client options.
func (c *ClientConfig) GetClientOptions() []define.BkApiClientOption {
	return c.ClientOptions
}
