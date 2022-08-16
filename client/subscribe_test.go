package client

import (
	"log"
	"testing"

	_ "github.com/ability-sh/abi-db/aws"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/runtime"
)

func TestSubscribe(t *testing.T) {

	config, err := runtime.GetConfigWithFile("../config.yaml")

	if err != nil {
		t.Fatal(err)
	}

	cfg := dynamic.GetWithKeys(config, []string{"services", "abi-mq"})

	s, err := NewSubscribe(dynamic.StringValue(dynamic.Get(cfg, "driver"), ""), cfg)

	if err != nil {
		t.Fatal(err)
	}

	s.Run(&SubscribeOptions{
		Topic: "test",
		Queue: "a",
		OnMessage: func(id string, body []byte) error {
			log.Println(id, string(body))
			return nil
		},
		OnWaitting: func() {
			log.Println("Waitting ...")
		},
	})

	select {}
}
