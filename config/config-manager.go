package config

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	viperObj "github.com/spf13/viper"
)

var configManager *ConfigurationManager

var once sync.Once

type ConfigurationManager struct {
	configuration *Configuration
}

func NewConfigurationManager(configPath string) (*ConfigurationManager, error) {

	var err error

	var config *Configuration

	once.Do(func() {
		config, err = createConfiguration(configPath)

		configManager = &ConfigurationManager{
			configuration: config,
		}
	})

	return configManager, err
}

func GetConfigManager() *ConfigurationManager {
	return configManager
}

func (cm *ConfigurationManager) GetDefaultGasLimit() *big.Int {
	return big.NewInt(cm.configuration.Gas.DefaultGasLimit)
}

func (cm *ConfigurationManager) GetDefaultGasLimit64() uint64 {
	return uint64(cm.configuration.Gas.DefaultGasLimit)
}

func (cm *ConfigurationManager) GetDefaultGasPrice() *big.Int {
	return big.NewInt(cm.configuration.Gas.DefaultGasPrice)
}

func (cm *ConfigurationManager) GetBlockDelay() *big.Int {
	return big.NewInt(cm.configuration.Network.BlockDelay)
}

func (cm *ConfigurationManager) GetBlockRetryInterval() time.Duration {
	return (time.Duration(cm.configuration.Network.BlockRetryInterval) * time.Second)
}

func (cm *ConfigurationManager) GetBlockRetryLimit() int {
	return cm.configuration.Network.BlockRetryLimit
}

func (cm *ConfigurationManager) GetTestEndPoint() string {
	return cm.configuration.Test.EndPoint
}

func (cm *ConfigurationManager) GetZeroAddress() common.Address {
	return common.HexToAddress(cm.configuration.Network.ZeroAddress)
}

func (cm *ConfigurationManager) GetTestChainID() uint8 {
	return cm.configuration.Test.TestChainID
}

func (cm *ConfigurationManager) GetTestRelayerThreshold() *big.Int {
	return big.NewInt(cm.configuration.Test.TestRelayerThreshold)
}

func (cm *ConfigurationManager) GetTestTimeout() time.Duration {
	return (time.Duration(cm.configuration.Test.TestTimeout )* time.Second)
}

func (cm *ConfigurationManager) GetExecuteBlockWatchLimit() int {
	return cm.configuration.Network.ExecuteBlockWatchLimit
}

func (cm *ConfigurationManager) GetTxRetryInterval() time.Duration {
	return time.Duration(cm.configuration.Network.TxRetryInterval) * time.Second
}

func (cm *ConfigurationManager) GetTxRetryLimit() int {
	return cm.configuration.Network.TxRetryLimit
}

func createConfiguration(configPath string) (*Configuration, error) {

	var configuration Configuration

	viper, err := bindDefaultConfigurations(viperObj.New(), configPath)

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

func bindDefaultConfigurations(viper *viperObj.Viper, configPath string) (*viperObj.Viper, error) {

	viper.SetConfigName("config")
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
