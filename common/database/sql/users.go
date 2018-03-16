package sql

import (
	"database/sql"

	"github.com/asciiu/gomo/common/models"
)

func GetUser(db *sql.DB, email string) (*models.User, error) {
	var u models.User
	err := db.QueryRow("SELECT id, first_name, last_name, email, email_verified, password_hash FROM users WHERE email = $1", email).
		Scan(&u.Id, &u.First, &u.Last, &u.Email, &u.EmailVerified, &u.PasswordHash)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
