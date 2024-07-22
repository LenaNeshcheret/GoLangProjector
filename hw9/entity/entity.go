package entity

import (
	"fmt"
	"sort"
	"sync"
)

type Student struct {
	ID     int
	Name   string
	Grades map[string]float64
}

type Class struct {
	ID       int
	Name     string
	Students []Student
	Teachers []Teacher
}
type Teacher struct {
	ID          int
	Name        string
	Classes     []Class
	TeacherCred TeacherCred
}
type TeacherCred struct {
	Username string
	Password string
}

type Storage struct {
	m           sync.Mutex
	lastID      int
	allClasses  map[int]Class
	allStudents map[int]Student
}

var MapAllClasses = make(map[int]Class)
var Classes []Class
var MapAllStudents = make(map[int]Student)

func NewStorage() *Storage {
	return &Storage{
		allClasses:  MapAllClasses,
		allStudents: MapAllStudents,
	}
}

func (s *Storage) GetAllClasses() []Class {
	s.m.Lock()
	defer s.m.Unlock()

	sort.Slice(Classes, func(i, j int) bool {
		return Classes[i].ID < Classes[j].ID
	})

	return Classes
}

func (s *Storage) GetClass(id int) Class {
	s.m.Lock()
	defer s.m.Unlock()

	class, exists := s.allClasses[id]
	if !exists {
		fmt.Printf("Class doesn't exist")
		return Class{}
	}
	return class
}

func (s *Storage) GetStudent(id int) Student {
	s.m.Lock()
	defer s.m.Unlock()

	student, exists := s.allStudents[id]
	if !exists {
		fmt.Printf("Student doesn't exist")
		return Student{}
	}
	return student
}
