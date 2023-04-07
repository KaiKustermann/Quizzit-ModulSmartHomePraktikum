# API

Contents:

- [Interfaces](#interfaces)
- [Specification](#specification)

## Interfaces

The backend offers a REST and a Websocket connection.

See below table to see which interface is used for each cause.

Name | REST/Socket | Content
--- | --- | --- 
Backend Available Check | REST | GET /health endpoint of backend (useful for client App to check if backend is online)
Profile CRUD | REST | GET profile(s), new profile, edit profile, delete profile
Open Websocket | Socket | Open Websocket connection
Gamestate Update from Backend | Socket | Send new gamestate
Action prompt from Backend | Socket | Send an action prompt to the active player (e.G. Roll Die or the question)
Decision from Player | Socket | Send user choice back to backend

## Specification

API Documentation is realized using:

Framework | Interface | Source | Comment
--- | --- | --- | ---
[Async API](https://www.asyncapi.com/docs/reference/specification/v2.6.0) | Websocket | [Source](./spec/websocket-asyncapi.yaml) | [Online Editor](https://studio.asyncapi.com/)
[Swagger](https://swagger.io/docs/specification) | REST | TODO | [Online Demo](https://petstore.swagger.io/), can load an own JSON file if available via http

Might find a proper way to provide shared objects to both asyncapi file and swagger file.
