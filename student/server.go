package student

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"example.com/service/grpc/student"
	"google.golang.org/grpc"
)

var _srv *grpc.Server

const PORT = "8089"

type studentServe struct {
	student.UnimplementedStudentServiceServer
}

func (s *studentServe) StudentOp(stream student.StudentService_StudentOpServer) error {
	for {
		e, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("event %v\n", e.Event)
		event, err := srs.Get(e.Event.TypeUrl)
		if err != nil {
			return fmt.Errorf("event[%v] error:%v", event, err)
		}

		ack, err := event.do(*e.Event)
		if err != nil {
			return fmt.Errorf("event[%v] error:%v", event, err)
		}
		stream.Send(ack)
	}
}

func (s *studentServe) SendHeartBeat(ctx context.Context, in *student.HeartBeat) (*student.Ack, error) {
	// cr.mu.RLock()
	// defer cr.mu.RUnlock()
	cr.mu.Lock()
	defer cr.mu.Unlock()
	fmt.Printf("heart beat %v\n", in.ClientId)
	cr.registry[in.ClientId] = time.Now()

	return &student.Ack{}, nil
}

func Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", PORT))
	if err != nil {
		fmt.Printf("failed to listen:%v\n", err)
		panic(err)
	}

	_srv = grpc.NewServer()

	student.RegisterStudentServiceServer(_srv, &studentServe{})
	go cr.offLine()
	if err := _srv.Serve(lis); err != nil {
		fmt.Printf("failed to start student service\n")
	}
}

func Stop() {
	fmt.Printf("stop student service")
}

type clientRegistry struct {
	registry map[string]time.Time
	mu       sync.RWMutex
}

var cr clientRegistry = clientRegistry{
	registry: make(map[string]time.Time, 4),
}

func (pct *clientRegistry) offLine() {
	for {
		for k, v := range pct.registry {
			if time.Since(v) > 9*time.Second {
				fmt.Printf("client %v off line\n", k)
				delete(pct.registry, k)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
