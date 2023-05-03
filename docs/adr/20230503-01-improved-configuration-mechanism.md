# Name
Improved Configuration Mechanism

## Status
Proposed

## Context
Our application needs a more robust and flexible configuration mechanism that can support the 12 factor app configuration guidelines, as well as various ways to configure the application and set values. Additionally, we need a hot reloading mechanism that can reload the configuration on the fly while the service is running.

## Decision
After considering several options, we have decided to use the open-source library Viper (https://github.com/spf13/viper) to implement our improved configuration mechanism. Viper supports a variety of config sources, including JSON, TOML, YAML, HCL, envfile, and Java properties files, as well as environment variables, remote config systems (such as etcd or Consul), and command line flags. Viper also supports hot reloading, which allows for changes to be made to the config files while the application is running.

## Consequences
By using Viper, we can simplify the process of reading and setting configuration values for our application, and ensure that our configuration follows the 12 factor app guidelines. Viper is a well-known and well-documented library that is widely used in the industry, which should make it easier for developers to understand and use the configuration mechanism. 

## Date
May 3, 2023

