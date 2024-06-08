package users

import (
	"api/internal/core"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type IMailService interface {
	Send(to, template string, content map[string]string) error
}

type users struct {
	db *sqlx.DB
	// mailService IMailService
}

func New(db *sqlx.DB) *users {
	return &users{db: db}
}

func (r *users) Add(user *core.User) (int, error) {
	var id int
	sql := "INSERT INTO users.users (fullname, email, password) VALUES ($1, $2, $3) RETURNING id;"
	err := r.db.QueryRowx(sql, user.Fullname, user.Email, user.Password).Scan(&id)
	return id, err
}

func (r *users) ExistsById(userId int) error {
	var id int
	query := "SELECT id FROM users.users WHERE id = $1 LIMIT 1"
	return r.db.QueryRow(query, userId).Scan(&id)
}

func (r *users) GetOneById(userId int) (*core.User, error) {
	var output core.User
	query := "SELECT * FROM users.users WHERE id = $1"
	err := r.db.QueryRowx(query, userId).StructScan(&output)
	return &output, err
}

func (r *users) GetByToken(token string) (*core.User, error) {
	var output core.User
	query := "SELECT * FROM users.users WHERE id = (SELECT user_id FROM users.sessions WHERE token = $1)"
	err := r.db.QueryRowx(query, token).StructScan(&output)
	return &output, err
}

func (r *users) GetByEmail(email string) (*core.User, error) {
	var output core.User
	sql := "SELECT * FROM users.users WHERE email = $1"
	err := r.db.QueryRowx(sql, email).StructScan(&output)
	return &output, err
}

func (r *users) EmailExists(email string) error {
	var userId int
	query := "SELECT id FROM users.users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&userId)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (r *users) CreateToken(userId int) (string, error) {
	token, err := core.GenerateRandomString(128)
	if err != nil {
		return "", nil
	}

	query := "INSERT INTO users.sessions (user_id, token) VALUES ($1, $2);"
	err = r.db.QueryRow(query, userId, token).Scan()
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return token, nil
}

func (r *users) Update(user *core.User) error {
	query := "UPDATE users.users SET fullname = $1, description = $2 WHERE id = $3"
	err := r.db.QueryRow(query, user.Fullname, user.Description, user.ID).Scan()
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (r *users) UpdateLastSeen(userId int) error {
	query := "UPDATE users.users SET last_seen = CURRENT_TIMESTAMP WHERE id = $1;"
	err := r.db.QueryRow(query, userId).Scan()
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (r *users) ChangePassword(userId int, password string) error {
	query := "UPDATE users.users SET password = $1 WHERE id = $2"
	return r.db.QueryRowx(query, password, userId).Err()
}

func (r *users) SendRecoveryKey(email string) error {
	user, err := r.GetByEmail(email)
	if err != nil {
		return err
	}

	key, err := core.GenerateRandomString(128)
	if err != nil {
		return err
	}

	expire := time.Now().Add(12 * time.Hour).UTC()

	query := "INSERT INTO users.recovery_keys (user_id, email, key, expire) VALUES ($1, $2, $3, $4);"
	if err := r.db.QueryRow(query, user.ID, email, key, expire).Scan(); err != nil && err != sql.ErrNoRows {
		return err
	}

	// mailContent := map[string]string{
	// 	"key": key,
	// }
	// if err := r.mailService.Send(ctx, email, "recovery", mailContent); err != nil {
	// 	return err
	// }

	return nil
}

func (r *users) DeleteRecoveryKey(key string) error {
	query := "DELETE FROM users.recovery_keys WHERE key = $1"
	err := r.db.QueryRow(query, key).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (r *users) GetUserIdByRecoveryKey(key string) (int, error) {
	var userId int
	query := "SELECT user_id FROM users.recovery_keys WHERE key = $1"
	err := r.db.QueryRow(query, key).Scan(&userId)
	return userId, err
}

func (r *users) SetPasswordByUserId(userId int, password string) error {
	query := "UPDATE users.users SET password = $1 WHERE id = $2"
	err := r.db.QueryRow(query, password, userId).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (r *users) UpdatePhoto(userId int, path string) error {
	query := "UPDATE users.users SET photo = $1 WHERE id = $2"
	err := r.db.QueryRow(query, path, userId).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}
