package students

type Storage interface {
	GetAllStudents() []Student
	CreateStudents(sts []Student) []Student
	UpdateStudent(id int, updatedStudent Student) bool
	UpdateStudentCourse(studentCourse StudentCourseProgress) StudentCourseProgress
	DeleteStudentByID(id int) bool
	SaveStudentCourseProgress(stCourse StudentCourseProgress) StudentCourseProgress
	SaveCourses(c []Course) []Course
	GetAllCourses() []Course
	GetCourseById(id int) Course
	UpdateCourseStatus(studentId int, courseId int, updateStatus UpdateStatus) StudentCourseProgress
}

type StudentService struct {
	s Storage
}

type CourseService struct {
	s   Storage
	sts StudentService
}

func NewStudentService(s Storage) *StudentService {
	return &StudentService{s: s}
}

func NewCourseService(s Storage, sts StudentService) *CourseService {
	return &CourseService{
		s:   s,
		sts: sts,
	}
}

func (s *StudentService) GetAllSt() []Student {
	return s.s.GetAllStudents()
}

func (s *StudentService) CreateSt(sts []Student) []Student {
	return s.s.CreateStudents(sts)
}

func (s *StudentService) UpdateSt(id int, updatedStudent Student) bool {
	return s.s.UpdateStudent(id, updatedStudent)
}

func (s *StudentService) UpdateStCourses(studentCourse StudentCourseProgress) StudentCourseProgress {
	return s.s.UpdateStudentCourse(studentCourse)
}

func (s *StudentService) DeleteStByID(id int) bool {
	return s.s.DeleteStudentByID(id)
}

func (s *StudentService) UpdateCourseStatus(studentId int, courseId int, updateStatus UpdateStatus) StudentCourseProgress {
	return s.s.UpdateCourseStatus(studentId, courseId, updateStatus)
}
func (s *CourseService) EnrollCourse(enrollCourse EnrollDTO) StudentCourseProgress {
	//
	courseId := enrollCourse.CourseId
	studentId := enrollCourse.StudentId
	course := s.s.GetCourseById(courseId)
	studentCourse := StudentCourseProgress{
		StudentId:   studentId,
		CourseId:    courseId,
		Name:        course.Name,
		Description: course.Description,
		Module:      course.Module,
		Status:      "enrolled",
	}
	//save
	courseProgress := s.sts.UpdateStCourses(studentCourse)
	return courseProgress
}

func (s *CourseService) CreateCourses(courses []Course) []Course {
	//save
	return s.s.SaveCourses(courses)
}
