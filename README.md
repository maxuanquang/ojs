# Online Judge System

## Table of contents

- [Description](#description)
  - [Key Features](#key-features)
- [Installation](#installation)
  - [Requirements](#requirements)
  - [Install](#install)
- [Usage](#usage)
  - [Running in local mode](#running-in-local-mode)
  - [Running in distributed mode](#running-in-distributed-mode)
- [Config](#config)
  - [CLI Arguments](#cli-arguments)
  - [Config files](#config-files)
- [Develop](#develop)
  - [Generate Code](#generate-code)

## Description

The Online Judge System (OJS) is a web-based platform designed for solving competitive programming problems with ease and efficiency. OJS can be setup easily for personal use and also posseses robust scalability to serve large-scale usage scenarios.

### Key Features

- Language Support: Offers support for a wide range of programming languages, providing flexibility for users to solve problems in their preferred language.
- Customization: Easy setup process allows customization for personal or organizational use.
- Scalability: Built with scalability in mind, capable of handling a high volume of users and submissions effortlessly.

## Installation

### Requirements

Docker or a Docker-compatible container runtime must be installed locally to run OJS process in local/distributed worker mode.

### Install

To install OJS, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/maxuanquang/ojs.git
```

2. Navigate to the project directory:

```bash
cd idm
```

3. Build project:

```bash
make build
```

4. Start all necessary services:

```bash
make docker-compose-prod-up
```

5. After all services have started up, we can start the project:

```bash
make run
```

By default, the HTTP server serves at `http:localhost:8081`

## Usage

### Running in local mode

To run OJS as a standalone server:

```bash
ojs standalone-server
```

### Running in distributed mode

To start HTTP host server:

```bash
ojs http-server
```

To start a worker process for evaluating submissions:

```bash
ojs worker
```

## Config

### CLI Arguments

| Argument | Description | Default Value |
| -------- | ----------- | ------------- |
| `--pull-image-at-startup`   | Whether to pull Docker images necessary for compiling and executing test case at startup. If set to true and Docker fails to pull any of the provided image, the program will exit with non-zero error code. | `true`                  |

### Config files

OJS use a YAML config file to configure its inner working in more details. By default, if no custom config file is provided, values in the config file [configs/local.yaml](configs/local.yaml) are used.

## Develop

### Generate code

OJS uses code generation for three purposes:

- **Compile-Time Dependency Injection**: Utilizes [github.com/google/wire](github.com/google/wire) for efficient compile-time dependency injection, ensuring robust and maintainable code.
- **gRPC Server, gRPC Gateway and OpenAPI Specifications Generation**: Automatically generates gRPC server and grpc gateway server to seamlessly handle RPC requests and HTTP requests, automatically produce OpenAPI specifications to promote interoperability and ease of integration.
- **JavaScript HTTP Client Generation**: Employs OpenAPI Generator to generate a JavaScript HTTP client, streamlining communication with the server-side components.

To execute all code generation processes, simply run:

```bash
make generate
```
