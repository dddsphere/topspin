# Name
Improved Configuration Mechanism

## Status
Proposed

## Context
In order to adhere to the 12 factor app configuration guidelines, our application requires a configuration mechanism that is both robust and flexible. It should offer multiple ways to configure the application and set values, and support hot reloading so that changes to the configuration can be made on the fly while the service is running. Furthermore, it's crucial that the application is compatible with standard Kubernetes mechanisms for configuration.

## Decision
After considering several options, we have decided to use the open-source library Viper (https://github.com/spf13/viper) to implement our improved configuration mechanism. Viper supports a variety of config sources, including JSON, TOML, YAML, HCL, envfile, and properties files, as well as environment variables, remote config systems (such as etcd or Consul), and command line flags. Viper also supports hot reloading, which allows for changes to be made to the config files while the application is running.

## Consequences
After evaluating several options, we have decided to implement our improved configuration mechanism using the open-source library Viper (https://github.com/spf13/viper). Viper offers a wide range of configuration sources, including JSON, TOML, YAML, HCL, envfile, and properties files, as well as environment variables, remote config systems (such as etcd or Consul), and command line flags. Additionally, Viper supports hot reloading, allowing for config files to be updated while the application is running. It's worth noting that Viper can also be configured to retrieve configuration data from traditional Kubernetes Config Maps and Secrets.

## Date
May 3, 2023

