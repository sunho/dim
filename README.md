# Introduction

Dim wraps echo to provide the dependecy injection for go web development.

It has been used to develop the server for [Minda](https://github.com/sdbx/minda), a game published in Steam.

# Features

## Simple service specific configuration

Each service is configured with service-specific yaml configuration file.

## Dependency injection

Dim can inject services to other services as well as echo specific structs such as middleware and route group.

# Getting Started

Try running the example application inside the `example` folder. You can build it by `go build` Pay attention to the stdout log and configuration files generated inside the working directory.