package students

import (
	"fmt"
	"sort"
	"sync"
)

type StudentCourse struct {
	studentId int
	courseId  int
}

type InMemStorage struct {
	m             sync.Mutex
	lastID        int
	students      map[int]Student
	courses       map[int]Course
	studentCourse map[StudentCourse]StudentCourseProgress
}

func NewInMemStorage() *InMemStorage {
	return &InMemStorage{
		students:      make(map[int]Student),
		courses:       make(map[int]Course),
		studentCourse: make(map[StudentCourse]StudentCourseProgress),
	}
}

func (s *InMemStorage) GetAllStudents() []Student {
	s.m.Lock()
	defer s.m.Unlock()

	var students = make([]Student, 0, len(s.students))

	for _, s := range s.students {
		students = append(students, s)
	}

	sort.Slice(students, func(i, j int) bool {
		return students[i].ID < students[j].ID
	})

	return students
}

func (s *InMemStorage) CreateStudents(sts []Student) []Student {
	s.m.Lock()
	defer s.m.Unlock()

	var resultStudents []Student
	fmt.Println("Trying to create students")
	for _, st := range sts {
		st.ID = s.lastID + 1
		s.students[st.ID] = st
		s.lastID++
		resultStudents = append(resultStudents, st)
		fmt.Printf("Created student with ID: %v\n", st.ID)
	}
	return resultStudents
}

func (s *InMemStorage) UpdateStudent(id int, updatedStudent Student) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, ok := s.students[id]
	if !ok {
		return false
	}

	updatedStudent.ID = id
	s.students[id] = updatedStudent
	return true
}

func (s *InMemStorage) UpdateStudentCourse(studentCourse StudentCourseProgress) StudentCourseProgress {
	s.m.Lock()
	defer s.m.Unlock()
	stId := studentCourse.StudentId
	// update courses in student
	st := s.students[stId]
	st.Courses = append(st.Courses, studentCourse)
	s.students[stId] = st
	return studentCourse
}

func (s *InMemStorage) DeleteStudentByID(id int) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, ok := s.students[id]
	if !ok {
		return false
	}

	delete(s.students, id)
	return true
}

func (s *InMemStorage) SaveStudentCourseProgress(stCourse StudentCourseProgress) StudentCourseProgress {
	s.m.Lock()
	defer s.m.Unlock()
	return stCourse
}

func (s *InMemStorage) SaveCourses(cs []Course) []Course {
	var result []Course
	for _, c := range cs {
		c.ID = s.lastID + 1
		s.courses[c.ID] = c
		s.lastID++
		result = append(result, c)
	}
	return result
}

func (s *InMemStorage) GetAllCourses() []Course {
	var courseRes []Course
	// Iterate over the map and append each value to the slice
	for _, course := range s.courses {
		courseRes = append(courseRes, course)
	}
	return courseRes
}

func (s *InMemStorage) GetCourseById(id int) Course {
	return s.courses[id]
}

func (s *InMemStorage) UpdateCourseStatus(studentId int, courseId int, updateStatus UpdateStatus) StudentCourseProgress {
	var scProgress StudentCourseProgress
	var index = -1
	for i, progress := range s.students[studentId].Courses {
		if progress.CourseId == courseId {
			scProgress = progress
			index = i
		}
	}
	totalProgress := float64(scProgress.Progress + updateStatus.Progress)
	if updateStatus.Status == "completed" {
		scProgress.Progress = 100
		scProgress.Status = "completed"
		scProgress.Mark = updateStatus.Mark
	} else if totalProgress >= 100 {
		scProgress.Mark =
			(updateStatus.Mark*float64(updateStatus.Progress) + scProgress.Mark*float64(scProgress.Progress)) / totalProgress
		scProgress.Progress = 100
		scProgress.Status = "completed"
	} else {
		scProgress.Mark = (updateStatus.Mark*float64(updateStatus.Progress) + scProgress.Mark*float64(scProgress.Progress)) / totalProgress
		scProgress.Progress += updateStatus.Progress
		scProgress.Status = "in progress"
	}
	if index != -1 {
		s.students[studentId].Courses[index] = scProgress
		return scProgress
	} // throw exception course id does not exist
	return StudentCourseProgress{}
}
