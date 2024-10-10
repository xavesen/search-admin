package models

import "encoding/json"

type User struct {
	Id         	string 		`json:"id,omitempty" bson:"_id,omitempty" validate:"omitempty,mongodb"`
	Login      	string 		`json:"login" validate:"required"`
	Password   	string 		`json:"password" validate:"required"`
	IndexLimit 	int    		`json:"index_limit" bson:"indexlimit" validate:"required"`
	Indexes		[]string	`json:"indexes,omitempty" validate:"omitempty"`
}

func (user *User) String() string {
	/*
	String representation of user struct is only used in logs.
	As logging passwords is a bad practice, password is censored.
	*/

	userNoPassword := *user
	userNoPassword.Password = "***"
	userJson, _ := json.Marshal(&userNoPassword)

	return string(userJson)
}