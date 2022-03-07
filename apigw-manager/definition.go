package manager

import (
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
		default:
			return nil, errors.Wrapf(ErrNotFound, "namespace: %s", namespace)
		}
	}

	return current, nil
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
