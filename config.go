package config

type ConfigReader interface {
	Get(key string) (interface{}, error)
	SetFromBytes(bytes []byte) error
}
