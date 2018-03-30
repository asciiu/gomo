package sql

import (
	"database/sql"

	"github.com/asciiu/gomo/user-service/models"
)

func DeleteUserHard(db *sql.DB, userId string) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", userId)
	return err
}

func DeleteUserSoft(db *sql.DB, userId string) error {
	_, err := db.Exec("UPDATE users SET deleted_on = now() WHERE id = $1", userId)
	return err
}

func FindUser(db *sql.DB, email string) (*models.User, error) {
	var u models.User
	err := db.QueryRow("SELECT id, first_name, last_name, email, email_verified, password_hash FROM users WHERE email = $1", email).
		Scan(&u.Id, &u.First, &u.Last, &u.Email, &u.EmailVerified, &u.PasswordHash)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func FindUserById(db *sql.DB, userId string) (*models.User, error) {
	var u models.User
	err := db.QueryRow("SELECT id, first_name, last_name, email, email_verified, password_hash FROM users WHERE id = $1", userId).
		Scan(&u.Id, &u.First, &u.Last, &u.Email, &u.EmailVerified, &u.PasswordHash)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func InsertUser(db *sql.DB, user *models.User) (*models.User, error) {
	sqlStatement := `insert into users (id, first_name, last_name, email, email_verified, password_hash) values ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(sqlStatement, user.Id, user.First, user.Last, user.Email, user.EmailVerified, user.PasswordHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserPassword(db *sql.DB, userId, hash string) error {
	_, err := db.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", hash, userId)
	return err
}

func UpdateUserInfo(db *sql.DB, user *models.User) (*models.User, error) {
	sqlStatement := `UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4`
	_, err := db.Exec(sqlStatement, user.First, user.Last, user.Email, user.Id)

	if err != nil {
		return nil, err
	}
	return user, nil
}
