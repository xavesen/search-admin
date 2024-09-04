package models

import "encoding/json"

type Filter struct {
	Id		string	`json:"id,omitempty" bson:"_id,omitempty" validate:"omitempty,mongodb"`
	Regex	string	`json:"regex" validate:"required"`
}

func (filter *Filter) String() string {
	filterJson, _ := json.Marshal(&filter)

	return string(filterJson)
}