# Configuration Management

The configuration package provides one source of truth for managing settings for infrastructure parameters 
such blockretryinterval, txretryinterval etc .
As well as provides bindings for environmental variables for easy access.


### Dependencies
- Viper https://github.com/spf13/viper

# Package content

The package consists of 3 files which include: 

- config.yaml  - All settings are defined here
- config.go    - provides a go data structure which maps to values defined in config.yaml
- config-manager -  binds config.yaml to config.go (as well as any environment variables) and provides getter functions and initializes a singleton of the configuration.

## Description 
- config.yaml 
  This YAML file contains the parameters to be consumed and their values. More values can be added here as the system grows

  config.yaml
  ```
    gas:
        defaultgaslimit: 6721975
        defaultgasprice: 20000000000
    network:
        blockdelay: 10
        blockretryinterval: 5
        blockretrylimit: 5
        zeroaddress: 0x0000000000000000000000000000000000000000
        executeblockwatchlimit: 100
        txretryinterval: 2
        txretrylimit: 10
  ```

- config.go Defines a struct that is mapped to the values defined in config.yaml 

  config.go

  ```
    package config

    type Gas struct {
        DefaultGasLimit int64
        DefaultGasPrice int64
    }

    type Network struct {
        BlockRetryInterval     int64
        BlockDelay             int64
        BlockRetryLimit        int
        ZeroAddress            string
        ExecuteBlockWatchLimit int
        TxRetryInterval        int64
        TxRetryLimit           int
    }

    type Configuration struct {
	    Gas     Gas
	    Network Network
    }

  ```

- configuration-manager.go 
  This file initializes the configaration by binding the values defined in config.yaml (as well as any specified environment variables) to the struct defined in config.go
  and instantiates a singleton of the configuration-manager which provides access to the bound properties via getter functions.

  ```
  package config


  type ConfigurationManager struct {
        configuration *Configuration
  }

  //initializes configuration manager - call this function at the application entry point to initialize the config
  func NewConfigurationManager(configPath string) (*ConfigurationManager, error) {

        // Perform initialization here by calling createConfiguration() defined below.
        // 1. bind config.yaml
        // 2. bind environment variables

        return configManager, err
   }
   
   // call this function any where in code to have access to the getter functions.
   func GetConfigManager() *ConfigurationManager {
	 return configManager
   }

   // getter function to return BlockRetryLimit 
   func (cm *ConfigurationManager) GetBlockRetryLimit() int {
      return cm.configuration.Network.BlockRetryLimit
   }

   // getter function to return TxRetryInterval
   func (cm *ConfigurationManager) GetTxRetryInterval() time.Duration {
	return time.Duration(cm.configuration.Network.TxRetryInterval) * time.Second
   }

   // getter function to return TxRetryLimit
   func (cm *ConfigurationManager) GetTxRetryLimit() int {
	 return cm.configuration.Network.TxRetryLimit
   }

   // getter function to return EXAMPLE_ENVIRONMENT_VARIABLE
   func (cm *ConfigurationManager) GetExampeEnvironmentVariable() string {
	return cm.configuration.EXAMPLE_ENVIRONMENT_VARIABLE
   }

   // called by NewConfigurationManager to bind config.yaml and environment variables to Configruration struct
   func createConfiguration(configPath string) (*Configuration, error) {

        var configuration Configuration

        //1. Call bindDefaultConfigurations
        //2. call bindEnvironmentVariables()
        //3  hydrate configuration with values from both bindings

	  return &configuration, nil
   }

   // reads config.yaml 
   func bindDefaultConfigurations(configPath string) (*viperObj.Viper, error) {

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

   // reads environment variables
   func bindEnvironmentVariables(viper *viperObj.Viper) {
	  viper.BindEnv("EXAMPLE_ENVIRONMENT_VARIABLE")
   }

   ```
