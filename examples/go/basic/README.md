# Basic Go Example

This example constructs basic Go types with no additional options.

The generated code was generated with the following command:

```sh
ecsgen generate --source-file "ecs_flat.yml" --output-plugin="gostruct" --opt-gostruct-package-name="main" --opt-gostruct-output-dir="."
```

To run the example, `cd` into the directory and run:

```sh
go run *.go | jq .
```

The result is:

```json
{
  "@timestamp": "2020-04-01T16:10:03.007692-07:00",
  "labels": {
    "foo": "bar"
  },
  "tags": [
    "production",
    "env2"
  ],
  "agent": {},
  "as": {
    "organization": {}
  },
  "client": {
    "as": {
      "organization": {}
    },
    "geo": {},
    "nat": {},
    "user": {
      "group": {}
    }
  },
  "cloud": {
    "account": {},
    "instance": {},
    "machine": {}
  },
  "code_signature": {},
  "container": {
    "image": {}
  },
  "destination": {
    "as": {
      "organization": {}
    },
    "geo": {},
    "nat": {},
    "user": {
      "group": {}
    }
  },
  "dll": {
    "code_signature": {},
    "hash": {},
    "pe": {}
  },
  "dns": {
    "question": {}
  },
  "ecs": {
    "version": "1.5.0"
  },
  "error": {},
  "event": {
    "created": "0001-01-01T00:00:00Z",
    "end": "0001-01-01T00:00:00Z",
    "ingested": "0001-01-01T00:00:00Z",
    "start": "0001-01-01T00:00:00Z"
  },
  "file": {
    "accessed": "0001-01-01T00:00:00Z",
    "code_signature": {},
    "created": "0001-01-01T00:00:00Z",
    "ctime": "0001-01-01T00:00:00Z",
    "hash": {},
    "mtime": "0001-01-01T00:00:00Z",
    "pe": {}
  },
  "geo": {},
  "group": {},
  "hash": {},
  "host": {
    "geo": {},
    "os": {},
    "user": {
      "group": {}
    }
  },
  "http": {
    "request": {
      "body": {}
    },
    "response": {
      "body": {}
    }
  },
  "interface": {},
  "log": {
    "origin": {
      "file": {}
    },
    "syslog": {
      "facility": {},
      "severity": {}
    }
  },
  "network": {
    "inner": {
      "vlan": {}
    },
    "vlan": {}
  },
  "observer": {
    "egress": {
      "interface": {},
      "vlan": {}
    },
    "geo": {},
    "ingress": {
      "interface": {},
      "vlan": {}
    },
    "os": {}
  },
  "organization": {},
  "os": {},
  "package": {
    "installed": "0001-01-01T00:00:00Z"
  },
  "pe": {},
  "process": {
    "code_signature": {},
    "hash": {},
    "parent": {
      "code_signature": {},
      "hash": {},
      "start": "0001-01-01T00:00:00Z",
      "thread": {}
    },
    "pe": {},
    "start": "0001-01-01T00:00:00Z",
    "thread": {}
  },
  "registry": {
    "data": {}
  },
  "related": {},
  "rule": {},
  "search": {
    "query": {},
    "results": {}
  },
  "server": {
    "as": {
      "organization": {}
    },
    "geo": {},
    "nat": {
      "ip": "192.168.2.4"
    },
    "user": {
      "group": {}
    }
  },
  "service": {
    "node": {}
  },
  "source": {
    "as": {
      "organization": {}
    },
    "geo": {},
    "nat": {},
    "user": {
      "group": {}
    }
  },
  "threat": {
    "tactic": {},
    "technique": {}
  },
  "tls": {
    "client": {
      "hash": {},
      "not_after": "0001-01-01T00:00:00Z",
      "not_before": "0001-01-01T00:00:00Z"
    },
    "server": {
      "hash": {},
      "not_after": "0001-01-01T00:00:00Z",
      "not_before": "0001-01-01T00:00:00Z"
    }
  },
  "trace": {},
  "transaction": {},
  "url": {},
  "user": {
    "group": {}
  },
  "user_agent": {
    "device": {},
    "os": {}
  },
  "vlan": {},
  "vulnerability": {
    "scanner": {},
    "score": {}
  }
}
```
