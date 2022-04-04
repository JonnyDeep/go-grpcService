package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/service/student"
)

func main() {
	fmt.Println("start student service")
	student.GetRegistry().Add("type.googleapis.com/student.AddStudentEvent", &student.AddStudentEvent{})
	student.GetRegistry().Add("type.googleapis.com/student.GetStudentEvent", &student.GetStudentEvent{})
	go student.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func(context.Context) {
		for {
			select {
			case <-ctx.Done():
				student.Stop()
				return
			default:
			}
		}
	}(ctx)
	time.Sleep(6 * time.Second)
}
