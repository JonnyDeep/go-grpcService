package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/service/grpc/student"
	studentService "example.com/service/student"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {

	ctx := context.Background()
	ip, port := "127.0.0.1", 8089
	client, err := studentService.NewStudentClient(ctx, ip, port)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	addEvent := &student.AddStudentEvent{
		Student: &student.Student{
			Id:   1000,
			Name: "Tom",
		},
	}
	event, err := anypb.New(addEvent)
	if err != nil {
		fmt.Printf("cannot marshal addStudnetEvent:%v\n", err)
	}
	reqest1 := &student.Event{
		ClientId: "10001",
		Event:    event,
	}
	client.Send(reqest1)

	getEvent := &student.GetStudentEvent{
		StudentId: 1000,
	}

	event, err = anypb.New(getEvent)
	if err != nil {
		fmt.Printf("cannot marshal getStudnetEvent:%v\n", err)
	}
	reqest2 := &student.Event{
		ClientId: "10001",
		Event:    event,
	}
	client.Send(reqest2)

	//heart beat
	go func(context.Context) {
		for {
			client.SendHeartBeat(ctx, &student.HeartBeat{
				ClientId: "10001",
			})
			time.Sleep(3 * time.Second)
		}
	}(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

}
