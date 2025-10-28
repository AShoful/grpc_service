package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" gorm:"unique" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignInInput struct {
	Username string `json:"username" gorm:"unique" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" gorm:"unique" validate:"required,min=4"`
	Author string `json:"author" validate:"required,min=4"`
	UserId uint   `json:"user_id"`
}

type UpdateBook struct {
	Title  *string `json:"title" gorm:"unique" validate:"omitempty,min=4"`
	Author *string `json:"author" validate:"omitempty,min=4"`
}
