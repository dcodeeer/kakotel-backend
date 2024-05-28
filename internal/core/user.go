package core

type User struct {
	ID          int     `json:"id" db:"id"`
	Email       string  `json:"email" db:"email"`
	Password    string  `json:"password" db:"password"`
	Phone       *string `json:"phone" db:"phone"`
	FirstName   *string `json:"firstname" db:"firstname"`
	LastName    *string `json:"lastname" db:"lastname"`
	Patronymic  *string `json:"patronymic" db:"patronymic"`
	DateOfBirth *string `json:"date_of_birth" db:"date_of_birth"`
	Photo       *string `json:"photo" db:"photo"`
	CreatedAt   string  `json:"created_at" db:"created_at"`
}
