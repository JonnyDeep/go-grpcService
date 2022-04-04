package student

import (
	"context"
	"fmt"
	"io"
	"time"

	"example.com/service/grpc/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StudentClient struct {
	client student.StudentServiceClient
	stream student.StudentService_StudentOpClient
}

func NewStudentClient(ctx context.Context, ip string, port int, opts ...grpc.CallOption) (*StudentClient, error) {
	address := fmt.Sprintf("%s:%d", ip, port)
	con, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("client can not connect server:%v", err)
	}

	client := student.NewStudentServiceClient(con)
	ss, err := client.StudentOp(ctx, opts...)
	if err != nil {
		fmt.Printf("%v.StudentOp(_)=_,%v", client, err)
	}
	sc := &StudentClient{
		client: student.NewStudentServiceClient(con),
		stream: ss,
	}
	go sc.recv(ctx)
	return sc, nil
}

func (s *StudentClient) SendHeartBeat(ctx context.Context, in *student.HeartBeat, opts ...grpc.CallOption) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err := s.client.SendHeartBeat(ctx, in, opts...)
	if err != nil {
		fmt.Printf("time out:%v\n", err)
	}

}

func (s *StudentClient) Send(event *student.Event) {
	s.stream.Send(event)
}

func (s *StudentClient) recv(ctx context.Context) {

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := s.stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				fmt.Printf("failed to receive data:%v\n", err)
			}
			fmt.Printf("receive data %v\n", in)
		}
	}()
	s.stream.CloseSend()
	<-waitc

}
