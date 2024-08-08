package students

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type service interface {
	GetAllSt() []Student
	CreateSt(st []Student) []Student
	UpdateSt(id int, updatedStudent Student) bool
	UpdateCourseStatus(studentId int, courseId int, updateStatus UpdateStatus) StudentCourseProgress
	DeleteStByID(id int) bool
}

type courseService interface {
	EnrollCourse(enrollCourse EnrollDTO) StudentCourseProgress
	CreateCourses(course []Course) []Course
}

type Handler struct {
	s  service
	cs courseService
}

func NewHandler(s service, cs courseService) Handler {
	return Handler{
		s:  s,
		cs: cs,
	}
}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	students := h.s.GetAllSt()

	err := json.NewEncoder(w).Encode(students)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreateStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student

	err := json.NewDecoder(r.Body).Decode(&students)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	students = h.s.CreateSt(students)

	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var updatedStudent Student

	err := json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idVal := r.PathValue("id")
	studentID, err := strconv.Atoi(idVal)
	if err != nil {
		fmt.Printf("Invalid id param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := h.s.UpdateSt(studentID, updatedStudent)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(updatedStudent)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	idVal := r.PathValue("id")

	studentID, err := strconv.Atoi(idVal)
	if err != nil {
		fmt.Printf("Invalid id param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := h.s.DeleteStByID(studentID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (h Handler) EnrollCourse(w http.ResponseWriter, r *http.Request) {
	idVal := r.PathValue("id")

	courseId, err := strconv.Atoi(idVal)
	if err != nil {
		fmt.Printf("Invalid id param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var enrollCourse EnrollDTO
	err = json.NewDecoder(r.Body).Decode(&enrollCourse)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	enrollCourse.CourseId = courseId
	course := h.cs.EnrollCourse(enrollCourse)

	err = json.NewEncoder(w).Encode(course)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (h Handler) CreateCourses(w http.ResponseWriter, r *http.Request) {

	var courses []Course
	err := json.NewDecoder(r.Body).Decode(&courses)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	returnCourses := h.cs.CreateCourses(courses)
	err = json.NewEncoder(w).Encode(&returnCourses)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h Handler) UpdateProgress(w http.ResponseWriter, r *http.Request) {

	studentIdVal := r.PathValue("id")
	studentId, err := strconv.Atoi(studentIdVal)
	if err != nil {
		fmt.Printf("Invalid id param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courseIdVal := r.PathValue("cid")
	courseId, err := strconv.Atoi(courseIdVal)
	if err != nil {
		fmt.Printf("Invalid id param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var updateStatus UpdateStatus

	err = json.NewDecoder(r.Body).Decode(&updateStatus)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//update student course progress
	status := h.s.UpdateCourseStatus(studentId, courseId, updateStatus)
	err = json.NewEncoder(w).Encode(&status)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
