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

package manager

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Definition represents a definition of a api gateway.
type Definition struct {
	definition map[string]interface{}
}

// Get sub definition.
func (d *Definition) Get(namespace string) (map[string]interface{}, error) {
	if namespace == "" {
		return d.definition, nil
	}

	current := d.definition
	for _, field := range strings.Split(namespace, ".") {
		if current == nil {
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}

		value, fond := current[field]
		if !fond {
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}

		switch realValue := value.(type) {
		case map[string]interface{}:
			current = realValue
		case map[interface{}]interface{}:
			// convert map[interface{}]interface{} to map[string]interface{}
			current = make(map[string]interface{})
			for k, v := range realValue {
				current[fmt.Sprintf("%v", k)] = v
			}
			return current, nil
		default:
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}
	}

	return current, nil
}

// GetArray Get sub array definition.
func (d *Definition) GetArray(namespace string) ([]map[string]interface{}, error) {
	current := d.definition
	for _, field := range strings.Split(namespace, ".") {
		if current == nil {
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}

		value, fond := current[field]
		if !fond {
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}

		switch realValue := value.(type) {
		case []interface{}:
			// convert []map[interface{}]interface{} to map[string]interface{}
			result := make([]map[string]interface{}, len(realValue))
			for i, v := range realValue {
				result[i] = v.(map[string]interface{})
			}
			return result, nil
		default:
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}
	}

	return []map[string]interface{}{}, nil
}

// NewDefinition creates a new definition from the given map.
func NewDefinition(definition map[string]interface{}) *Definition {
	return &Definition{
		definition: definition,
	}
}

// NewDefinitionFromYaml unmarshal the given yaml string to a definition.
func NewDefinitionFromYaml(content []byte) (*Definition, error) {
	var definition map[string]interface{}
	err := yaml.Unmarshal(content, &definition)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}

	return NewDefinition(definition), nil
}
