package sql

import (
	"database/sql"
	"time"

	"github.com/asciiu/gomo/api/models"
)

func FindRefreshToken(db *sql.DB, selector string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := db.QueryRow("SELECT id, user_id, selector, token_hash, expires_on FROM refresh_tokens WHERE selector = $1", selector).
		Scan(&t.ID, &t.UserID, &t.Selector, &t.TokenHash, &t.ExpiresOn)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func InsertRefreshToken(db *sql.DB, token *models.RefreshToken) (*models.RefreshToken, error) {
	sqlStatement := `insert into refresh_tokens (id, user_id, selector, token_hash, expires_on) values ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, token.ID, token.UserID, token.Selector, token.TokenHash, token.ExpiresOn)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func DeleteRefreshToken(db *sql.DB, token *models.RefreshToken) (*models.RefreshToken, error) {
	sqlStatement := `delete from refresh_tokens where id = $1`
	_, err := db.Exec(sqlStatement, token.ID)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func DeleteRefreshTokenBySelector(db *sql.DB, selector string) error {
	sqlStatement := `delete from refresh_tokens where selector = $1`
	_, err := db.Exec(sqlStatement, selector)
	return err
}

func DeleteStaleTokens(db *sql.DB, expiresOn time.Time) error {
	sqlStatement := `delete from refresh_tokens where expires_on < $1`
	_, err := db.Exec(sqlStatement, expiresOn)
	return err
}

func UpdateRefreshToken(db *sql.DB, token *models.RefreshToken) (*models.RefreshToken, error) {
	sqlStatement := `update refresh_tokens set selector = $1, token_hash = $2, expires_on = $3 where user_id = $4 and id = $5`
	_, err := db.Exec(sqlStatement, token.Selector, token.TokenHash, token.ExpiresOn, token.UserID, token.ID)
	if err != nil {
		return nil, err
	}
	return token, nil
}
