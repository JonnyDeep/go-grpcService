package student

import (
	"fmt"

	"example.com/service/grpc/student"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/types/known/anypb"
)

type GRPCMESSAGE interface {
	student.AddStudentEvent | student.GetStudentEvent
}

type event interface {
	do(any.Any) (*student.Ack, error)
}

type AddStudentEvent struct{}

func (a *AddStudentEvent) do(any any.Any) (*student.Ack, error) {
	stu := &student.AddStudentEvent{}
	err := any.UnmarshalTo(stu)
	if err != nil {
		fmt.Printf("message format error:%v\n", err)
	}
	if err = stuDb.addStudent(stu.Student); err != nil {
		fmt.Println(err)
	}
	return &student.Ack{}, nil
}

type GetStudentEvent struct{}

func (a *GetStudentEvent) do(any any.Any) (*student.Ack, error) {
	event := &student.GetStudentEvent{}
	err := any.UnmarshalTo(event)
	if err != nil {
		fmt.Printf("message format error:%v\n", err)
	}
	stu, err := stuDb.getStudentDetail(event.GetStudentId())
	if err != nil {
		fmt.Println(err)
	}
	ack, err := anypb.New(stu)

	if err != nil {
		return &student.Ack{}, err
	}
	return &student.Ack{
		Ack: ack,
	}, nil
}

type EventRegistry struct {
	registry map[string]event
}

var srs = &EventRegistry{
	registry: make(map[string]event, 4),
}

func GetRegistry() *EventRegistry {
	return srs
}

func (sr *EventRegistry) Add(key string, val event) {
	sr.registry[key] = val
}

func (sr *EventRegistry) Get(key string) (event, error) {
	e, ok := sr.registry[key]
	if !ok {
		return nil, fmt.Errorf("no %v event", key)
	}
	return e, nil
}
