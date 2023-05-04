# Name
Improved Configuration Mechanism

## Status
Proposed

## Context
Our application needs a more robust and flexible configuration mechanism that can support the 12 factor app configuration guidelines, as well as various ways to configure the application and set values. Additionally, we need a hot reloading mechanism that can update the configuration on the fly while the service is running. It is also important that the application can be also configured using standard Kubernetes mechanisms such as ConfigMaps and Secrets in the case our application needs to run on a cluster.

## Decision
After considering several options, we have decided to use the open-source library Viper (https://github.com/spf13/viper) to implement our improved configuration mechanism. Viper supports a variety of config sources, including JSON, TOML, YAML, HCL, envfile, and Java properties files, as well as environment variables, remote config systems (such as etcd or Consul), and command line flags. Viper also supports hot reloading, which allows for changes to be made to the config files while the application is running.

We plan to implement a custom Viper provider that will handle configuration when running outside of Kubernetes. This provider will allow us to use a local YAML file for configuration when the application is running outside of Kubernetes, and the Kubernetes-based provider when running inside a Kubernetes cluster. This will provide us with a flexible and consistent configuration mechanism across different deployment environments.

## Consequences
By using Viper, we can simplify the process of reading and setting configuration values for our application, and ensure that our configuration follows the 12 factor app guidelines. Viper is a well-known and well-documented library that is widely used in the industry, which should make it easier for developers to understand and use the configuration mechanism. 

## Date
May 3, 2023

