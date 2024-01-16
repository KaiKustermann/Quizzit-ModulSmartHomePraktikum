# Backend Server

Server Code, the heart of the system.

# Index

- [Backend Server](#backend-server)
- [Index](#index)
- [Scope / Purpose](#scope--purpose)
- [Development](#development)
  - [Requirements](#requirements)
  - [Setup](#setup)
  - [Run](#run)
  - [Debug in VS-Code](#debug-in-vs-code)
  - [Debug Connectivity](#debug-connectivity)
  - [Test](#test)
  - [Build](#build)
  - [Code format](#code-format)
  - [Testing via Websocket](#testing-via-websocket)
- [Production](#production)
  - [Create deployment](#create-deployment)
  - [System Requirements:](#system-requirements)
- [Related Content](#related-content)

# Scope / Purpose

Backend has the following tasks:

* Handles Game State
* Connects to the [hybrid die](https://gitlab.mi.hdm-stuttgart.de/quizzit/hybrid-die)
* future: Connects to the Hue Lights

The backend is in flow-control when the game is running.

The api and communication flow with the client web app is described in [the api folder](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/tree/main/api).

# Development

Directory structure as seen in ["Golang Standard Project Layout"](https://github.com/golang-standards/project-layout)

## Requirements

* Golang 1.20 or higher
* npm to generate GO code from api spec

## Setup

Download dependencies and generate GO code from api spec.

    go mod download
    npm i
    go generate cmd/quizzit/quizzit.go (generates golang types from the ./api/websocket-asyncapi.yaml)

## Run

To start the server:
    
    go run cmd/quizzit/quizzit.go

The server will start listening on http://localhost:8080.

**If you are using Windows and want to use the [Hybrid Die](../hybrid-die/), make sure you are in the same `private` Network as the die. `Private` referes to your system's security settings of this Network. For windows, see the screenshot below.**
![Private network settings](./assets/img/private-network-windows.png)

**If you want to work on the hybrid die itself, [this trick](#configuration-to-simplify-hybrid-die-development) may be useful**

## Debug in VS-Code

Possible VS-Code `.vscode/launch.json` configuration:

```json
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Quizzit GO",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/cmd/quizzit/quizzit.go",
      "cwd": "${workspaceFolder}/backend",
      "args": [
        "-log-level", "debug"
      ],
      "trace": "trace"
    }
  ]
}
```

## Debug Connectivity

A useful tool to debug connectivity issues is Wireshark.

When running on a RaspberryPi we recommend `tshark`, a CLI version of wireshark

```sh
# Installing tshark
sudo apt-get update
sudo apt-get upgrade
# During install, when prompted if non-root users need to capture packages you do not need to select 'yes' (select 'no')
sudo apt-get install tshark

# Capture packages on port 7778 and display the content as text
sudo tshark -i wlan0 -o data.show_as_text:TRUE -T fields -f "port 7778" -e data.text
```

## Test

    go test ./test

## Build

To build the binary:

    go build cmd/quizzit/quizzit.go

This will create a binary named quizzit in the project directory.

## Code format

Golang has a built-in command-line tool called go fmt that automatically formats Go source code. The go fmt command formats your code according to a set of rules defined in the Go code [style guidelines](https://go.dev/doc/effective_go#formatting).

## Testing via Websocket

The best way to test, is using the corresponding [frontend](https://gitlab.mi.hdm-stuttgart.de/quizzit/client-web-app).

However to test new states and messages you can use this lovely page here: https://websocketking.com/ and connect to this backend - probably via: `ws://localhost:8080/ws`.

Any message requires at least the `messageType`. 
Refer to the [api specification](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/blob/main/api/websocket-asyncapi.yaml), the active gamestate or the gameloop that is printed in the beginning.

Example for an often-used Message - A simple confirm:

```json
{
  "messageType": "player/generic/Confirm"
}
```
# Production

## Create deployment

To create a deployable gitlab release, go to the [tag section](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/tags/) of the backend repository in gitlab. 

Create a new tag with a specified new version number (e.g. v0.0.1-RC1).

A new release will be created, which can be found in the [deploy section](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/releases).

In production the backend runs as binary on a RaspberryPi 4B along with the client web app. For more information regarding deploying to production, look into [the RaspberryPi Installation Guide](https://gitlab.mi.hdm-stuttgart.de/quizzit/raspberry-pi/-/blob/main/Installation-Guide-RaspberryPi.md).

## System Requirements:

- tested on RaspberryPi 4B
- should run on most systems with:
  * Windows or Linux as operating system
  * Wifi card



# Related Content

* [API Specification](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/tree/main/api)
* [Backend Wiki](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/home)
* [Game Loop](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/home)
* [RaspberryPi Installation Guide](https://gitlab.mi.hdm-stuttgart.de/quizzit/raspberry-pi/-/blob/main/Installation-Guide-RaspberryPi.md)
* [Hybrid Die](https://gitlab.mi.hdm-stuttgart.de/quizzit/hybrid-die)
* [Go Generate](https://go.dev/blog/generate)
