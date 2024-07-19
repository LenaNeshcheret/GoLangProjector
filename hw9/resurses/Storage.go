package resursestype

import (
	"GoLangProjector/hw9/entity"
)

type Storage interface {
	GetAllClasses() []entity.Class
	GetClass(id int) entity.Class
	GetStudent(id int) entity.Student
}
