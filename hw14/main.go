package main

import (
	"GoLangProjector/hw14/students"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	storage := students.NewInMemStorage()

	studentService := students.NewStudentService(storage)
	courseService := students.NewCourseService(storage, *studentService)

	handler := students.NewHandler(studentService, courseService)

	mux.HandleFunc("GET /students", handler.GetAll)
	mux.HandleFunc("POST /students", handler.CreateStudents)
	mux.HandleFunc("PUT /students/{id}", handler.UpdateStudent)
	mux.HandleFunc("DELETE /students/{id}", handler.DeleteStudent)
	mux.HandleFunc("PUT /students/{id}/course/{cid}", handler.UpdateProgress)
	mux.HandleFunc("PATCH /courses/{id}/enroll", handler.EnrollCourse)
	mux.HandleFunc("POST /courses", handler.CreateCourses)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error is occurred: ", err.Error())
		return
	}
}
