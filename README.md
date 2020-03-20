# ecsgen

**THIS IS A WORK IN PROGRESS**

This project aims to fix some of the shortcomings of the [Elastic Common Schema](https://github.com/elastic/ecs) and it's representation as Go types. Elastic has shipped a tool that generates Go types from the ECS YAML definitions, but it has shortcomings outlined below. This project aims to implement it better.

## Install

```sh
go install github.com/gen0cide/ecsgen/cmd/ecsgen
```

## Todo

[] Generate actual Go code
[] Consider strategies to minimize nil pointer deferencing
[] Resolve field types
[] Perform field validation based on type
[] Test it? (how?)
[] Generate documentation in a nice way, especially for the types that we had to create
[] Find someone to bother at Elastic who can help make sense of some of their design decisions

Meh
