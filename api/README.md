# API Definitions

- [API Definitions](#api-definitions)
  - [Communication in game (Websockets)](#communication-in-game-websockets)
  - [Communications Settings (HTTP/REST)](#communications-settings-httprest)
    - [Swagger UI](#swagger-ui)
- [Futher Reading](#futher-reading)

The API definitions are used to generate the interface-describing code for both client and backend.

## Communication in game (Websockets)

For communication between the game backend and the client web app, websockets are used.   

The API and communication flow with the client web app is described in the [websocket-asyncapi.yaml](./websocket-asyncapi.yaml). You can use the website [asyncapi](https://studio.asyncapi.com/) to visualize the api document (copy and paste it in).

A detailed overview of the gameloop and communication via websockets is also available in the [backend WIKI](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/home).

## Communications Settings (HTTP/REST)

Before the game starts we need no realtime duplex communication, so we can use REST.

The API ist described in [openapi.yaml](./openapi.yaml). You can use the [swagger editor](https://editor-next.swagger.io/) to visualize the api document (copy and paste it in).

### Swagger UI

Windows:

    docker run --rm -p 8000:8080 -e SWAGGER_JSON=/local/api/openapi.yaml -v %CD%:/local swaggerapi/swagger-ui

Linux:

    docker run --rm -p 8000:8080 -e SWAGGER_JSON=/local/api/openapi.yaml -v $PWD:/local swaggerapi/swagger-ui

# Futher Reading

* [See the API WIKI](https://gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/-/wikis/Updating-the-API).
* https://github.com/OpenAPITools/openapi-generator/blob/master/docs/generators/go.md
* https://github.com/OpenAPITools/openapi-generator/blob/master/docs/generators/typescript-rxjs.md
