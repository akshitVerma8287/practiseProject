package models

type Student struct {
	ID int `json:"id" gorm:"Id primary key"`
	Name string `json:"name" gorm:"Name"`
	Age int `json:"age" gorm:"Age"`
	Email string `json:"email" gorm:"Email"`
}

type Admin struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}