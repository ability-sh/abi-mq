package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/ability-sh/abi-db/source"
)

type Subscribe struct {
	s source.Source
}

func NewSubscribe(driver string, config interface{}) (*Subscribe, error) {
	s, err := source.NewSource(driver, config)
	if err != nil {
		return nil, err
	}
	return &Subscribe{s: s}, nil
}

type SubscribeOptions struct {
	Topic      string
	Queue      string
	OnMessage  func(id string, body []byte) error
	OnError    func(err error)
	OnWaitting func()
}

var re_id, _ = regexp.Compile(`([a-zA-Z0-9]+)\.msg$`)

func (s *Subscribe) Run(opt *SubscribeOptions) func() {

	c, cancel := context.WithCancel(context.Background())

	go func() {

		prefix := fmt.Sprintf("%s/%s/", opt.Topic, opt.Queue)

		tk := time.NewTicker(time.Second * 6)

		defer tk.Stop()

		running := true

		for running {

			rs, err := s.s.Query(prefix, "/")

			if err != nil {
				if opt.OnError != nil {
					opt.OnError(err)
				} else {
					log.Println(err)
				}
			} else {

				for {

					key, err := rs.Next()

					if err != nil {
						if err == io.EOF {
							break
						} else if opt.OnError != nil {
							opt.OnError(err)
						} else {
							log.Println(err)
						}
						break
					}

					rr := re_id.FindStringSubmatch(key)

					if len(rr) > 1 {
						id := rr[1]
						body, err := s.s.Get(key)
						if err != nil {
							if opt.OnError != nil {
								opt.OnError(err)
							} else {
								log.Println(err)
							}
							break
						}
						err = opt.OnMessage(id, body)
						if err != nil {
							if opt.OnError != nil {
								opt.OnError(err)
							} else {
								log.Println(err)
							}
							break
						}
						err = s.s.Del(key)
						if err != nil {
							if opt.OnError != nil {
								opt.OnError(err)
							} else {
								log.Println(err)
							}
							break
						}
					}

				}
				rs.Close()
			}

			if opt.OnWaitting != nil {
				opt.OnWaitting()
			}
			select {
			case <-tk.C:
			case <-c.Done():
				running = false
			}
		}

	}()

	return cancel

}
