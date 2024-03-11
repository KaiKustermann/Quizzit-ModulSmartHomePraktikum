# Quizzit Configuration

- [Quizzit Configuration](#quizzit-configuration)
  - [Configuration Sources](#configuration-sources)
  - [Configuration Precedence](#configuration-precedence)

## Configuration Sources

Configuration can be supplied by various sources:

1. FLAGs when starting the application
2. The (main) configuration file (for the path, run program with `--help`)
3. Configuration through the user. This also creates files in the `userConfigDir` (for the path, run program with `--help`)

## Configuration Precedence

The program has a default config that ensures it can run without extra config.
[See here](./configuration.go)

This configuration is the base, that will get patched by the other configuration sources as follows:

1. Patches through the (main) config file
2. Patches through FLAGs
3. Patches through the user config (files)

As a user changes the configuration through the API, these changes will also always overrule any other configuration.
Of course they will be preserved to disk into the mentioned files in the `userConfigDir` (for the path, run program with `--help`)
