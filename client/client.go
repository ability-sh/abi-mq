package client

import (
	"context"

	"github.com/ability-sh/abi-lib/errors"
	"github.com/ability-sh/abi-mq/pb"
	"google.golang.org/grpc"
)

type Topic struct {
	cli   pb.ServiceClient
	topic string
}

func (t *Topic) Send(c context.Context, body string) (string, error) {
	rs, err := t.cli.Send(c, &pb.SendTask{Topic: t.topic, Body: body})
	if err != nil {
		return "", err
	}
	if rs.Errno == 200 {
		return rs.Id, nil
	}
	return "", errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

func (t *Topic) QueueCreate(c context.Context, queue string) error {
	rs, err := t.cli.QueueCreate(c, &pb.QueueCreateTask{Topic: t.topic, Queue: queue})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

func (t *Topic) QueueRemove(c context.Context, queue string) error {
	rs, err := t.cli.QueueRemove(c, &pb.QueueRemoveTask{Topic: t.topic, Queue: queue})
	if err != nil {
		return err
	}
	if rs.Errno == 200 {
		return nil
	}
	return errors.Errorf(rs.Errno, "%s", rs.Errmsg)
}

type Client struct {
	cli pb.ServiceClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{cli: pb.NewServiceClient(conn)}
}

func (c *Client) Topic(topic string) *Topic {
	return &Topic{topic: topic, cli: c.cli}
}
