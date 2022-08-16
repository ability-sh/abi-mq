package srv

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ability-sh/abi-micro/micro"
	"github.com/ability-sh/abi-mq/pb"
	"google.golang.org/grpc"
)

type server struct {
}

func (srv *server) Send(c context.Context, task *pb.SendTask) (*pb.SendResult, error) {

	if task.Topic == "" {
		return &pb.SendResult{Errno: 400, Errmsg: "input topic error"}, nil
	}

	ctx := micro.GetContext(c)

	ss, err := GetMQSource(ctx, SERVICE_ABI_MQ)

	if err != nil {
		return &pb.SendResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	prefix := fmt.Sprintf("%s/", task.Topic)

	rs, err := ss.Query(prefix, "/")

	if err != nil {
		return &pb.SendResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	defer rs.Close()

	id := strconv.FormatInt(ctx.Runtime().NewID(), 36)
	data := []byte(task.Body)

	for {

		key, err := rs.Next()

		if err != nil {
			if err == io.EOF {
				break
			}
			return &pb.SendResult{Errno: 500, Errmsg: err.Error()}, nil
		}

		if strings.HasSuffix(key, "/") {
			err = ss.Put(fmt.Sprintf("%s%s.msg", key, id), data)
			if err != nil {
				return &pb.SendResult{Errno: 500, Errmsg: err.Error()}, nil
			}
		}
	}

	return &pb.SendResult{Errno: 200, Id: id}, nil

}

func (srv *server) QueueCreate(c context.Context, task *pb.QueueCreateTask) (*pb.QueueCreateResult, error) {

	if task.Topic == "" {
		return &pb.QueueCreateResult{Errno: 400, Errmsg: "input topic error"}, nil
	}

	if task.Queue == "" {
		return &pb.QueueCreateResult{Errno: 400, Errmsg: "input queue error"}, nil
	}

	ctx := micro.GetContext(c)

	ss, err := GetMQSource(ctx, SERVICE_ABI_MQ)

	if err != nil {
		return &pb.QueueCreateResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	err = ss.Put(fmt.Sprintf("%s/%s/metafile.json", task.Topic, task.Queue), []byte("{}"))

	if err != nil {
		return &pb.QueueCreateResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	return &pb.QueueCreateResult{Errno: 200}, nil
}

func (srv *server) QueueRemove(c context.Context, task *pb.QueueRemoveTask) (*pb.QueueRemoveResult, error) {

	if task.Topic == "" {
		return &pb.QueueRemoveResult{Errno: 400, Errmsg: "input topic error"}, nil
	}

	if task.Queue == "" {
		return &pb.QueueRemoveResult{Errno: 400, Errmsg: "input queue error"}, nil
	}

	ctx := micro.GetContext(c)

	ss, err := GetMQSource(ctx, SERVICE_ABI_MQ)

	if err != nil {
		return &pb.QueueRemoveResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	prefix := fmt.Sprintf("%s/", task.Topic)

	rs, err := ss.Query(prefix, "")

	if err != nil {
		return &pb.QueueRemoveResult{Errno: 500, Errmsg: err.Error()}, nil
	}

	defer rs.Close()

	for {

		key, err := rs.Next()

		ctx.Println(key)

		if err != nil {
			if err == io.EOF {
				break
			}
			return &pb.QueueRemoveResult{Errno: 500, Errmsg: err.Error()}, nil
		}

		err = ss.Del(key)

		if err != nil {
			return &pb.QueueRemoveResult{Errno: 500, Errmsg: err.Error()}, nil
		}

	}

	return &pb.QueueRemoveResult{Errno: 200}, nil
}

func Reg(s *grpc.Server) {
	pb.RegisterServiceServer(s, &server{})
}
