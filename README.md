# ecsgen

**THIS IS A WORK IN PROGRESS**

This project aims to fix some of the shortcomings of the [Elastic Common Schema](https://github.com/elastic/ecs) and it's representation as Go types. More description to be written as the tool evolves.

## Install

```sh
go install github.com/gen0cide/ecsgen/cmd/ecsgen
```

## Todo

- [ ] Generate actual Go code
- [ ] Consider strategies to minimize nil pointer deferencing
- [ ] Resolve field types
- [ ] Perform field validation based on type
- [ ] Test it? (how?)
- [ ] Generate documentation in a nice way, especially for the types that we had to create
- [ ] Find someone to bother at Elastic who can help make sense of some of their design decisions

Meh
