package converter

import (
	"GoLangProjector/hw6/hw9/dto"
	"GoLangProjector/hw6/hw9/entity"
)

type Converter interface {
	Convert(interface{}) interface{}
}
type ClassConverter struct{}

func (c *ClassConverter) Convert(input interface{}) interface{} {
	classes, ok := input.([]entity.Class)
	if !ok {
		return nil
	}
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

func (c *StudentConverter) Convert(input interface{}) interface{} {
	student, ok := input.(entity.Student)
	if !ok {
		return nil
	}
	dtoStudent := dto.Student{
		ID:     student.ID,
		Name:   student.Name,
		Grades: student.Grades,
	}

	return dtoStudent
}
