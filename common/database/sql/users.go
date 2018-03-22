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

func InsertUser(db *sql.DB, user *models.User) (*models.User, error) {
	sqlStatement := `insert into users (id, first_name, last_name, email, email_verified, password_hash, salt) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(sqlStatement, user.Id, user.First, user.Last, user.Email, user.EmailVerified, user.PasswordHash, user.Salt)

	if err != nil {
		return nil, err
	}
	return user, nil
}
