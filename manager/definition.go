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

// Definition represents a definition of an api gateway.
type Definition struct {
	definition map[string]any
}

// Get gets the sub definition under the namespace, such as "apigateway" or "a.b".
func (d *Definition) Get(namespace string) (map[string]any, error) {
	if namespace == "" {
		return d.definition, nil
	}

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
		case map[string]any:
			current = realValue
		case map[any]any:
			// convert map[any]any to map[string]any
			current = make(map[string]any)
			for k, v := range realValue {
				current[fmt.Sprintf("%v", k)] = v
			}
		default:
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}
	}

	return current, nil
}

// GetArray gets the array of objects under the namespace, such as "stages" or "a.b".
func (d *Definition) GetArray(namespace string) ([]map[string]any, error) {
	parent, field := "", namespace
	if i := strings.LastIndex(namespace, "."); i >= 0 {
		parent, field = namespace[:i], namespace[i+1:]
	}
	current, err := d.Get(parent)
	if err != nil {
		return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
	}
	value, found := current[field]
	if !found {
		return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
	}
	items, ok := value.([]any)
	if !ok {
		return nil, errors.Wrapf(errors.New("not supported type"), "namespace: %s", namespace)
	}
	// convert []any to []map[string]any
	result := make([]map[string]any, len(items))
	for i, v := range items {
		item, ok := v.(map[string]any)
		if !ok {
			return nil, errors.Wrapf(errors.New("not supported type"), "namespace: %s", namespace)
		}
		result[i] = item
	}
	return result, nil
}

// NewDefinition creates a new definition from the given map.
func NewDefinition(definition map[string]any) *Definition {
	return &Definition{
		definition: definition,
	}
}

// NewDefinitionFromYaml unmarshals the given yaml to a definition.
func NewDefinitionFromYaml(content []byte) (*Definition, error) {
	var definition map[string]any
	err := yaml.Unmarshal(content, &definition)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}

	return NewDefinition(definition), nil
}
