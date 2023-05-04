# Backend Server

Server Code, the heart of the system.

# Index

* [Scope / Purpose](#scope--purpose)
* [Development](#development)
    * [Requirements](#requirements)
    * [Setup](#setup)
    * [Run](#run)
    * [Test](#test)
    * [Build](#build)
    * [Code Formatting](#code-format)
* [Related Content](#related-content)

# Scope / Purpose

Backend has the following tasks:

* Handles Game State
* Handles (Player) Profiles
* future: Connects to the https://gitlab.mi.hdm-stuttgart.de/quizzit/hybrid-die ?
* future: Connects to the Hue Lights

The backend is in flow-control when the game is running.
For configuration phase the client is in control using REST.

# Development

Directory structure as seen in ["Golang Standard Project Layout"](https://github.com/golang-standards/project-layout)

## Requirements

* Golang 1.20 or higher
* npm to generate GO code from api spec

## Setup

Download dependencies and generate GO code from api spec.

    go mod download
    npm i
    go generate cmd/quizzit/quizzit.go

## Run

To start the server:
    
    go run cmd/quizzit/quizzit.go

The server will start listening on http://localhost:8080.

## Test

    go test ./test

## Build

To build the binary:

    go build cmd/quizzit/quizzit.go

This will create a binary named quizzit in the project directory.

## Code format

Golang has a built-in command-line tool called go fmt that automatically formats Go source code. The go fmt command formats your code according to a set of rules defined in the Go code [style guidelines](https://go.dev/doc/effective_go#formatting).

*Could be implemented in CI/CD in the future.*

## Testing Websocket

Can use this lovely page here: https://websocketking.com/ and connect to `ws://localhost:8080/ws`.

### Submit an Answer

```json
{
  "CorrelationId": "a-cor-id",
  "MessageType": 0,
  "Body": {
    "QuestionId": "question-0",
    "AnswerId": "C"
  }
}
```

# Related Content

* [API Specification](./api)
* [Backend Wiki](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/home)
* [Game Loop](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/home)
* [Go Generate](https://go.dev/blog/generate)
