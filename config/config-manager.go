package config

import (
	"fmt"
	"time"

	viperObj "github.com/spf13/viper"
)

type ConfigurationManager struct {
	configuration *Configuration
}

func NewConfigurationManager(configPath string, fileName string) (*ConfigurationManager, error) {

	var err error

	var config *Configuration

	config, err = createConfiguration(configPath, fileName)

	configManager := &ConfigurationManager{
		configuration: config,
	}

	return configManager, err
}

func (cm *ConfigurationManager) GetDefaultGasLimit() int {
	return cm.configuration.Gas.DefaultGasLimit
}

func (cm *ConfigurationManager) GetDefaultGasPrice() int {
	return cm.configuration.Gas.DefaultGasPrice
}

func (cm *ConfigurationManager) GetDefaultGasPrice() int {
	return (cm.configuration.Network.BlockRetryInterval * time.Second)
}

func createConfiguration(configPath string, fileName string) (*Configuration, error) {

	var configuration Configuration

	viper, err := bindDefaultConfigurations(viperObj.New(), configPath, fileName)

	if err != nil {
		return nil, err
	}

	bindEnvironmentVariables(viper)

	err = viper.Unmarshal(&configuration)

	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func bindDefaultConfigurations(viper *viperObj.Viper, configPath string, fileName string) (*viperObj.Viper, error) {

	viper.SetConfigName(fileName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("..config error", err)
		return nil, err
	}

	return viper, nil
}

func bindEnvironmentVariables(viper *viperObj.Viper) {
	//viper.BindEnv("<ENVIRONMENT_VARIABLE_NAME>")
}
