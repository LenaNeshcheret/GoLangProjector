package dto

type Class struct {
	ID       int
	Name     string
	Students []string `json:"studentNames"`
}

type Student struct {
	ID     int
	Name   string
	Grades map[string]float64
}
