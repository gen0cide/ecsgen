# ecsgen

This project aims to create a well defined way to perform code generation on [Elastic Common Schema](https://github.com/elastic/ecs) definitions. Originally, this was intended to be specific to Go, but is built in a way that implementing new output generators is easy.

## Install

```sh
go install github.com/gen0cide/ecsgen/cmd/ecsgen
```

## Usage

To use `ecsgen`, there are a few options:

```
--source-file value                   Path to the generated ecs_flat.yml file containing ECS definitions. [$ECSGEN_SOURCE_FILE]
--whitelist value                     Regular expression that denotes which ECS keys to allow into the model. (Can be used multiple times). [$ECSGEN_WHITELIST_VALUE]
--blacklist value                     Regular expression that denotes which ECS keys to explicitly forbid into the model. (Can be used multiple times). [$ECSGEN_BLACKLIST_VALUE]
--output-plugin value                 Enable an output generator plugin. Can be used multiple times. Possible values: debug, gostruct [$ECSGEN_OUTPUT_PLUGIN]
```

The only required ones are `--source-file` that points to the ecs_flat.yml ECS definition, as well as at least one `--output-plugin`.

## Examples

Check out the examples/ folder.

## Current Output Plugins

The current list of usable output plugins is:

### `gostruct`

Gostruct is used to generate Go code for an ECS object. It has a few options:

```
--opt-gostruct-package-name value     Name of the Go package for the generated code. [$ECSGEN_OPT_GOSTRUCT_PACKAGE_NAME]
--opt-gostruct-output-dir value       Path to the directory where the generated code should be written. [$ECSGEN_OPT_GOSTRUCT_OUTPUT_DIR]
--opt-gostruct-output-filename value  Destination filename for the generated code. (default: generated_ecs.go) [$ECSGEN_OPT_GOSTRUCT_OUTPUT_FILENAME]
--opt-gostruct-marshal-json           Include a json.Marshaler implementation that removes empty fields. (default: false) [$ECSGEN_OPT_GOSTRUCT_MARSHAL_JSON]
```

The `--opt-gostruct-marshal-json` is shown in the examples/go/with-json-marshaling example directory.
