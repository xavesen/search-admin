package models

type User struct {
	Id 			int		`json:"id"`
	Login		string	`json:"login"`
	Password	string	`json:"password"`
	IndexLimit	int		`json:"index_limit"`
}