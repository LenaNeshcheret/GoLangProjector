package converter

import (
	"GoLangProjector/hw9/dto"
	"GoLangProjector/hw9/entity"
)

type ClassConverter struct{}

func (c *ClassConverter) Convert(classes []entity.Class) []dto.Class {
	dtoClasses := make([]dto.Class, len(classes))
	for i, class := range classes {
		dtoStudents := make([]string, len(class.Students))
		for j, student := range class.Students {
			dtoStudents[j] = student.Name
		}
		dtoClasses[i] = dto.Class{
			ID:       class.ID,
			Name:     class.Name,
			Students: dtoStudents,
		}
	}
	return dtoClasses
}

type StudentConverter struct{}

func (c *StudentConverter) Convert(student entity.Student) dto.Student {
	dtoStudent := dto.Student{
		ID:     student.ID,
		Name:   student.Name,
		Grades: student.Grades,
	}

	return dtoStudent
}
