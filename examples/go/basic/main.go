package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	evt := Base{}
	evt.Server.NAT.IP = "192.168.2.4"
	evt.ECS.Version = "1.5.0"
	evt.AtTimestamp = time.Now()
	if evt.Labels == nil {
		evt.Labels = map[string]interface{}{}
	}
	evt.Labels["foo"] = "bar"
	evt.Tags = append(evt.Tags, "production", "env2")

	blob, err := json.Marshal(&evt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(blob))
}
