package models

import "github.com/kamva/mgm/v3"


type User struct {
	mgm.DefaultModel `bson:",inline"`

	Name string `json:"name" bson:"name" validate:"required"`
	Email string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type Arithmetic struct {
	mgm.DefaultModel `bson:".inline"`

	Equation string `json:"equation" bson:"equation" validate:"required"`
	Action string `json:"action" bson:"action" validate:"required"`
	User string `json:"user" bson:"user" validate:"required"`
}
