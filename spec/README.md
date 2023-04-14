# Specification

API Documentation and generating source code from it.

# Index

* [Tools](#tools)
* [Publish Models](#publish-models)
    * [Typescript](#publish-models-for-typescript)
* [Related Content](#related-content)

# Tools

Framework | Interface | Source | Comment
--- | --- | --- | ---
[Async API](https://www.asyncapi.com/docs/reference/specification/v2.6.0) | Websocket | [Source](./spec/websocket-asyncapi.yaml) | [Online Editor](https://studio.asyncapi.com/)
[Swagger](https://swagger.io/docs/specification) | REST | [Source](./spec/swagger.yaml) | [Online Editor](https://editor.swagger.io/)

While it is generally possible to import definitions from another file, 
[the asyncapi-studio does not support file imports yet](https://github.com/asyncapi/studio/pull/538)
and swagger-editor also does not show the schemata nicely.

Therefore for now the specifications contain duplicate definitions.
One for async-api and one for openapi.

In the future we might find a proper way to provide shared objects to both asyncapi file and swagger file.

# Publish Models

In order to generate and publish the models after changing the API you must initially follow these steps:

1. [Create/Get an access token](https://im-gitlab.hdm-stuttgart.de/groups/quizzit/-/settings/access_tokens)
with role `developer` and scope `api`
2. `cp example.env .env` and use your (new) token as `GITLAB_TOKEN`.

After this setup you can now run the commands for the given languages...

## Publish Models for Typescript

Set an appropriate version in [package.json](./generate/typescript/package.json).

    docker-compose run --rm -it generate-typescript

*If you received a 403 when publishing it is either an invalid token or maybe the version already exists!*

# Related Content

* [Gitlab NPM Registry Docs](https://docs.gitlab.com/ee/user/packages/npm_registry/)
