package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type Yaml struct {
	config map[string]interface{}
	lock   sync.RWMutex
}

func NewYaml() *Yaml {
	return &Yaml{}
}

func (c *Yaml) Get(key string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	skeys := strings.Split(key, ".")
	return getPathValue(skeys, c.config)
}

func (c *Yaml) SetFromBytes(data []byte) error {
	var rawConfig interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return err
	}

	untypedConfig, ok := rawConfig.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("config is not a map")
	}

	config, err := convertKeysToStrings(untypedConfig)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.config = config
	return nil
}

var ErrNotFound = errors.New("can not find key in config")

func getPathValue(keys []string, cfg interface{}) (interface{}, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("param keys must not be empty: got %v", keys)
	}
	var (
		c  map[string]interface{}
		ok bool
		v  interface{}
	)
	if c, ok = cfg.(map[string]interface{}); !ok {
		return nil, ErrNotFound
	}
	if v, ok = c[keys[0]]; ok {
		if len(keys) == 1 {
			return v, nil
		}
		return getPathValue(keys[1:], v)
	}
	return nil, ErrNotFound
}

func convertKeysToStrings(m map[interface{}]interface{}) (map[string]interface{}, error) {
	n := make(map[string]interface{})

	for k, v := range m {
		str, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("config key is not a string")
		}

		if vMap, ok := v.(map[interface{}]interface{}); ok {
			var err error
			v, err = convertKeysToStrings(vMap)
			if err != nil {
				return nil, err
			}
		}

		n[str] = v
	}

	return n, nil
}
