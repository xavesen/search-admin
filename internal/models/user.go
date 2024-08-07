package models

import "encoding/json"

type User struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	IndexLimit int    `json:"index_limit" bson:"indexlimit"`
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