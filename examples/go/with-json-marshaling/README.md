# JSON Marshaler Go Example

This example constructs basic Go types with the json.Marshaler implementation that removes empty fields from the resulting JSON.

The generated code was generated with the following command:

```sh
ecsgen generate --source-file "ecs_flat.yml" --output-plugin="gostruct" --opt-gostruct-package-name="main" --opt-gostruct-output-dir="." --opt-gostruct-marshal-json
```

To run the example, `cd` into the directory and run:

```sh
go run *.go | jq .
```

The result is:

```json
{
  "@timestamp": "2020-04-01T16:10:33.039405-07:00",
  "ecs": {
    "version": "1.5.0"
  },
  "labels": {
    "foo": "bar"
  },
  "server": {
    "nat": {
      "ip": "192.168.2.4"
    }
  },
  "tags": [
    "production",
    "env2"
  ]
}
```
