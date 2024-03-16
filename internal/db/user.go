package db

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/models"
)

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	result := models.User{}
	queryOrder := `
	SELECT email, password, role FROM users WHERE email = $1
	`
	err := db.dbConnect.QueryRow(ctx, queryOrder, email).Scan(&result.Email, &result.Password, &result.Role)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
