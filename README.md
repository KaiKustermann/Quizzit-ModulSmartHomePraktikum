# Backend Server

Server Code, the heart of the system.

# Index

* [Scope / Purpose](#scope--purpose)
* [Related Content](#related-content)

# Scope / Purpose

Backend has the following tasks:

* Handles Game State
* Handles (Player) Profiles
* future: Connects to the https://gitlab.mi.hdm-stuttgart.de/quizzit/hybrid-die ?
* future: Connects to the Hue Lights

The backend is in flow-control when the game is running.
For configuration phase the client is in control using REST.

# Related Content

* [API Specification](./spec)
* [Backend Wiki](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/home)
* [Game Loop](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/home)


# Prerequisites
Golang 1.20 or higher

# Setup
## Clone the repository:
git clone https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server.git

## Install dependencies:
```go mod download```

## Run
To start the server, run: <br>
```go run main.go```<br>
The server will start listening on http://localhost:8080.

## Build
To build the binary, run: <br>
```go build -o quizzit-server main.go``` <br>
This will create a binary named quizzit-server in the project directory.

## Testing
To run all tests, run: <br>
```go test```

# Code format
Golang has a built-in command-line tool called go fmt that automatically formats Go source code. The go fmt command formats your code according to a set of rules defined in the Go code [style guidelines](https://go.dev/doc/effective_go#formatting). 
Could be implemented in CI/CD in the future.
