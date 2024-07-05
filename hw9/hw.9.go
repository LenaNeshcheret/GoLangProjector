package main

import (
	"GoLangProjector/hw6/hw9/converter"
	"GoLangProjector/hw6/hw9/entity"
	resursestype "GoLangProjector/hw6/hw9/resurses"
	"fmt"
	"net/http"
	"strconv"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/classes", classResource.GetAllClasses)
	mux.HandleFunc("/classes/{id}", classResource.GetClass)
	mux.HandleFunc("/students/{id}", checkAuth(studentResource.GetUserById))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error is occurred: ", err.Error())
		return
	}
}

func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		idUser := r.PathValue("id")
		id, err := strconv.Atoi(idUser)
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
		}
		var teachersList []entity.Teacher
		classes := classResource.S.GetAllClasses()
		for _, class := range classes {
			students := class.Students
			for _, student := range students {
				if student.ID == id {
					teachersList = append(teachersList, class.Teachers...)
				}
			}
		}
		var isAuthorised bool
		for _, teacher := range teachersList {
			if teacher.TeacherCred.Username == username && teacher.TeacherCred.Password == password {
				isAuthorised = true
			}
		}
		if !isAuthorised {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

var teachers []entity.Teacher
var classResource resursestype.ClassResource
var studentResource resursestype.StudentResource

func init() {

	students := []entity.Student{
		{ID: 1, Name: "John Doe", Grades: map[string]float64{"Math": 90, "Science": 85}},
		{ID: 2, Name: "Jane Smith", Grades: map[string]float64{"Math": 92, "Science": 88}},
		{ID: 3, Name: "Emily Johnson", Grades: map[string]float64{"Math": 85, "Science": 89}},
		{ID: 4, Name: "Michael Brown", Grades: map[string]float64{"Math": 87, "Science": 84}},
		{ID: 5, Name: "Sarah Davis", Grades: map[string]float64{"Math": 93, "Science": 91}},
		{ID: 6, Name: "David Wilson", Grades: map[string]float64{"Math": 78, "Science": 80}},
		{ID: 7, Name: "Laura Martinez", Grades: map[string]float64{"Math": 88, "Science": 86}},
		{ID: 8, Name: "James Garcia", Grades: map[string]float64{"Math": 91, "Science": 87}},
		{ID: 9, Name: "Sophia Martinez", Grades: map[string]float64{"Math": 90, "Science": 92}},
		{ID: 10, Name: "Christopher Lee", Grades: map[string]float64{"Math": 84, "Science": 83}},
	}

	entity.Classes = []entity.Class{
		{
			ID:   2,
			Name: "Chemistry 101",
			Students: []entity.Student{
				students[5], students[6], students[7], students[8], students[9],
			},
		},
		{
			ID:   1,
			Name: "Physics 101",
			Students: []entity.Student{
				students[0], students[1], students[2], students[3], students[4],
			},
		},
		{
			ID:   3,
			Name: "Biology 101",
			Students: []entity.Student{
				students[0], students[2], students[4], students[6], students[8],
			},
		},
		{
			ID:   4,
			Name: "Math 101",
			Students: []entity.Student{
				students[1], students[3], students[5], students[7], students[9],
			},
		},
		{
			ID:   5,
			Name: "History 101",
			Students: []entity.Student{
				students[0], students[1], students[3], students[5], students[7],
			},
		},
	}

	for _, c := range entity.Classes {
		entity.MapAllClasses[c.ID] = c
	}

	for _, s := range students {
		entity.MapAllStudents[s.ID] = s
	}
	storage := entity.NewStorage()

	classResource = resursestype.ClassResource{
		S: storage,
		C: &converter.ClassConverter{},
	}
	studentResource = resursestype.StudentResource{
		S: storage,
		C: &converter.StudentConverter{},
	}

	teachers = []entity.Teacher{
		{ID: 1, Name: "Mr. Smith",
			TeacherCred: entity.TeacherCred{
				Username: "Olena",
				Password: "1111",
			}},
		{ID: 2, Name: "Ms. Johnson",
			TeacherCred: entity.TeacherCred{
				Username: "Olena",
				Password: "2222",
			}},
		{ID: 3, Name: "Dr. Brown",
			TeacherCred: entity.TeacherCred{
				Username: "Olena",
				Password: "3333",
			}},
	}

	teachers[0].Classes = []entity.Class{entity.Classes[0], entity.Classes[1]}
	teachers[1].Classes = []entity.Class{entity.Classes[2], entity.Classes[3]}
	teachers[2].Classes = []entity.Class{entity.Classes[4]}

	entity.Classes[0].Teachers = []entity.Teacher{teachers[0]}
	entity.Classes[1].Teachers = []entity.Teacher{teachers[0]}
	entity.Classes[2].Teachers = []entity.Teacher{teachers[1]}
	entity.Classes[3].Teachers = []entity.Teacher{teachers[1]}
	entity.Classes[4].Teachers = []entity.Teacher{teachers[2]}

}

//func printClassAndTeacherDetails() {
//	//Printing class details
//	for _, class := range entity.Classes {
//		fmt.Printf("Class ID: %d, Name: %s\n", class.ID, class.Name)
//		fmt.Println("Students:")
//		for _, student := range class.Students {
//			fmt.Printf("\tStudent ID: %d, Name: %s\n", student.ID, student.Name)
//		}
//		fmt.Println("Teachers:")
//		for _, teacher := range class.Teachers {
//			fmt.Printf("\tTeacher ID: %d, Name: %s\n", teacher.ID, teacher.Name)
//		}
//		fmt.Println()
//	}
//
//	// Printing teacher details
//	for _, teacher := range teachers {
//		fmt.Printf("Teacher ID: %d, Name: %s\n", teacher.ID, teacher.Name)
//		fmt.Println("Classes:")
//		for _, class := range teacher.Classes {
//			fmt.Printf("\tClass ID: %d, Name: %s\n", class.ID, class.Name)
//		}
//		fmt.Println()
//	}
//}
