package service

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("abi-mq", func(name string, config interface{}) (micro.Service, error) {
		return newClientService(name, config), nil
	})
}
