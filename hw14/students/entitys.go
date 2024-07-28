package students

type Student struct {
	ID      int
	Name    string
	Email   string
	Rating  float32
	Courses []StudentCourseProgress
}
type Course struct {
	ID          int
	Name        string
	Description string
	Module      []Module
}
type EnrollDTO struct {
	ID        int
	CourseId  int `json:"courseId"`
	StudentId int `json:"studentId"`
}

type Module struct {
	ID   int
	Name string
}

type UpdateStatus struct {
	Status   interface{}
	Progress int
	Mark     float64
}

type StudentCourseProgress struct {
	ID          int
	StudentId   int
	CourseId    int
	Name        string
	Description string
	Module      []Module
	Status      string
	Progress    int
	Mark        float64
}
