package comic

import (
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

// comicViper defines the methods of Viper used by Comic
type comicViper interface {
	SetConfigName(in string)
	AddConfigPath(in string)
	AutomaticEnv()
	SetEnvKeyReplacer(r *strings.Replacer)
	ReadInConfig() error
	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
	IsSet(key string) bool
	AllKeys() []string
}

// mockViper is a Viper stand-in for Comic testing
type mockViper struct {
	cfg  interface{}
	keys map[string]bool
}

func (m *mockViper) SetConfigName(in string) {}

func (m *mockViper) AddConfigPath(in string) {}

func (m *mockViper) AutomaticEnv() {}

func (m *mockViper) SetEnvKeyReplacer(r *strings.Replacer) {}

func (m *mockViper) ReadInConfig() error {
	return nil
}

func (m *mockViper) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) (err error) {
	if m.cfg == nil {
		return
	}

	reflect.ValueOf(rawVal).Elem().Set(reflect.ValueOf(m.cfg).Elem())

	return
}

func (m *mockViper) IsSet(key string) bool {
	return m.keys[key]
}

func (m *mockViper) AllKeys() []string {
	allKeys := make([]string, len(m.keys))

	for key := range m.keys {
		allKeys = append(allKeys, key)
	}

	sort.Strings(allKeys)

	return allKeys
}
