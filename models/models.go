package models

type Student struct {
	Name  string  `json:"name"`
	Class int     `json:"class"`
	Grade float64 `json:"grade"`
}
type StudentList struct {
	Student
	Info string `json:"info"`
}
