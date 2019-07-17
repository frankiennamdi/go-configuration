package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/frankiennamdi/go-configuration/configuration"
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
	err := configuration.Bind(mappings, &config)
	fmt.Printf("%+v\n", config)
	require.NoError(t, err)
	assert.Equal(t, config.Server.Host, "localhost")
	assert.Equal(t, config.Server.Port, "9002")
}

func TestBindConfig(t *testing.T) {
	mappings := map[string]string{
		"server_host": "localhost",
		"server_port": "9002",
	}
	config := Config{}
	err := configuration.InitializeConfig(mappings, &config)
	fmt.Printf("%+v\n", config)
	require.NoError(t, err)
	assert.Equal(t, config.Server.Host, "localhost")
	assert.Equal(t, config.Server.Port, "9002")
}

func TestConfigWithBindEnvironment(t *testing.T) {
	os.Setenv("server_host", "localhost")
	os.Setenv("server_port", "9002")
	defer os.Unsetenv("server_host")
	defer os.Unsetenv("server_port")
	config := Config{}
	err := configuration.BindEnvironment(&config)
	fmt.Printf("%+v\n", config)
	require.NoError(t, err)
	assert.Equal(t, config.Server.Host, "localhost")
	assert.Equal(t, config.Server.Port, "9002")
}
