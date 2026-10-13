/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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
	yaml "gopkg.in/yaml.v3"
)

// Definition represents a definition of a api gateway.
type Definition struct {
	definition map[string]any
}

// Get sub definition.
func (d *Definition) Get(namespace string) (map[string]any, error) {
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
		case map[string]any:
			current = realValue
		case map[any]any:
			// convert map[any]any to map[string]any
			current = make(map[string]any)
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
func (d *Definition) GetArray(namespace string) ([]map[string]any, error) {
	current := d.definition
	for _, field := range strings.Split(namespace, ".") {
		if current == nil {
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}
		value, found := current[field]
		if !found {
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}
		switch realValue := value.(type) {
		case []any:
			// convert []map[any]any to map[string]any
			result := make([]map[string]any, len(realValue))
			for i, v := range realValue {
				result[i] = v.(map[string]any)
			}
			return result, nil
		default:
			return nil, errors.Wrapf(errors.New("not supported type"), "namespace: %s", namespace)
		}
	}

	return []map[string]any{}, nil
}

// NewDefinition creates a new definition from the given map.
func NewDefinition(definition map[string]any) *Definition {
	return &Definition{
		definition: definition,
	}
}

// NewDefinitionFromYaml unmarshal the given yaml string to a definition.
func NewDefinitionFromYaml(content []byte) (*Definition, error) {
	var definition map[string]any
	err := yaml.Unmarshal(content, &definition)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}

	return NewDefinition(definition), nil
}
