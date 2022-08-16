package srv

import (
	"fmt"

	"github.com/ability-sh/abi-db/source"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/micro"
)

const (
	SERVICE_ABI_MQ = "abi-mq"
)

type MQService struct {
	config interface{}   `json:"-"`
	name   string        `json:"-"`
	Driver string        `json:"driver"`
	Source source.Source `json:"-"`
}

func newMQService(name string, config interface{}) *MQService {
	return &MQService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *MQService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *MQService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *MQService) OnInit(ctx micro.Context) error {

	dynamic.SetValue(s, s.config)

	ss, err := source.NewSource(s.Driver, s.config)

	if err != nil {
		return err
	}

	s.Source = ss

	return nil
}

/**
* 校验服务是否可用
**/
func (s *MQService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *MQService) Recycle() {

}

func GetMQService(ctx micro.Context, name string) (*MQService, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(*MQService)
	if ok {
		return ss, nil
	}
	return nil, fmt.Errorf("service %s not instanceof *MQService", name)
}

func GetMQSource(ctx micro.Context, name string) (source.Source, error) {
	s, err := GetMQService(ctx, name)
	if err != nil {
		return nil, err
	}
	return s.Source, nil
}
