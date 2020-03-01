package configuration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Server struct {
	Host string `config:"host"`
	Port string `config:"port"`
}

type Config struct {
	Server Server `config:"server"`
}

func TestExpandMap(t *testing.T) {
	mappings := map[string]string{
		"server_host": "localhost",
		"server_port": "9002",
	}
	config := Config{}
	binder := New()
	err := binder.Bind(mappings, &config)
	fmt.Printf("%+v\n", config)
	require.NoError(t, err)
	assert.Equal(t, config.Server.Host, "localhost")
	assert.Equal(t, config.Server.Port, "9002")
}

func TestBindConfigFromYaml_With_No_Defaults(t *testing.T) {
	setEnv("SERVER_HOST", "localhost")
	setEnv("SERVER_PORT", "9002")
	defer unsetEnv("SERVER_HOST")
	defer unsetEnv("SERVER_PORT")
	data, err := ioutil.ReadFile("config_with_no_default.yml")
	require.NoError(t, err)
	config := Config{}
	binder := New()
	initErr := binder.InitializeConfigFromYaml(data, &config)
	require.NoError(t, initErr)
	fmt.Printf("%+v\n", config)
	assert.Equal(t, config.Server.Host, "localhost")
	assert.Equal(t, config.Server.Port, "9002")
}

func TestBindConfigFromYaml_With_No_Defaults_Fails_When_Env_Missing(t *testing.T) {

	data, err := ioutil.ReadFile("config_with_no_default.yml")
	require.NoError(t, err)
	config := Config{}
	binder := New()
	initErr := binder.InitializeConfigFromYaml(data, &config)
	require.Error(t, initErr)
}

func TestBindConfigFromYaml_With_Defaults(t *testing.T) {
	data, err := ioutil.ReadFile("config_with_default.yml")
	require.NoError(t, err)
	config := Config{}
	binder := New()
	initErr := binder.InitializeConfigFromYaml(data, &config)
	require.NoError(t, initErr)
	fmt.Printf("%+v\n", config)
	assert.Equal(t, config.Server.Host, "127.0.0.1")
	assert.Equal(t, config.Server.Port, "5000")
}

func unsetEnv(key string) {
	if err := os.Unsetenv(key); err != nil {
		log.Fatalf("%+v", err)
	}
}

func setEnv(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Fatalf("%+v", err)
	}
}
