package core

type User struct {
	ID          int     `json:"id" db:"id"`
	Email       string  `json:"email" db:"email"`
	Password    string  `json:"password" db:"password"`
	Phone       *string `json:"phone" db:"phone"`
	Fullname    *string `json:"fullname" db:"fullname"`
	Description *string `json:"description" db:"description"`
	DateOfBirth *string `json:"date_of_birth" db:"date_of_birth"`
	LastSeen    string  `json:"last_seen" db:"last_seen"`
	Photo       *string `json:"photo" db:"photo"`
	CreatedAt   string  `json:"created_at" db:"created_at"`
}
