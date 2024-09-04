package models

type Filter struct {
	Id		string	`json:"id,omitempty" bson:"_id,omitempty" validate:"omitempty,mongodb"`
	Regex	string	`json:"regex" validate:"required"`
}