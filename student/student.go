package student

import (
	"fmt"
	"sync"

	"example.com/service/grpc/student"
)

func init() {
	stuDb = &StudentDb{
		stus: make(Students, 1),
	}

	stuDb.stus = append(stuDb.stus, student.Student{
		Id:   10000,
		Name: "tom",
	}, student.Student{
		Id:   10001,
		Name: "jerry",
	})
}

type StudentDb struct {
	stus Students
	mu   sync.RWMutex
}

var stuDb *StudentDb

type Students []student.Student

func (sdb *StudentDb) getStudentDetail(id int64) (*student.Student, error) {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()
	for _, v := range sdb.stus {
		if v.Id == id {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("can not find id %d student", id)
}

func (sdb *StudentDb) addStudent(stu *student.Student) error {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()
	for _, v := range sdb.stus {
		if v.Id == stu.Id {
			return fmt.Errorf("the student[id:%d] is exist", v.Id)
		}
	}
	sdb.stus = append(sdb.stus, *stu)
	return nil
}
