package database

import (
	"database/sql"
	"fmt"

	"github.com/ValeriyL01/balance-service/internal/models"
)

type UserDB struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{db: db}
}

func (u UserDB) CreateUser(user *models.User) error {

	query := `
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at
`
	err := u.db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}
	return nil
}
func (u UserDB) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
              FROM users WHERE username = $1`

	user := &models.User{}
	err := u.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("юзер не найден")
	}
	return user, err
}
func (u *UserDB) GetUserByID(id int64) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
              FROM users WHERE id = $1`

	user := &models.User{}
	err := u.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("юзер не найден")
	}
	return user, err
}
