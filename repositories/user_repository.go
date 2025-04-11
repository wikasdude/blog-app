package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wikasdude/blog-backend/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	UpdateUser(user *model.User) error
	IsEmailTaken(email string, excludeUserID int) (bool, error)
	DeleteUser(id int) error
	GetUserByEmail(email string) (*model.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	var dbName string
	err := r.db.QueryRow("SELECT current_database()").Scan(&dbName)
	if err != nil {
		log.Fatal("❌ Failed to fetch current database name:", err)
	}
	fmt.Println("✅ Connected to database:", dbName)
	query := `INSERT INTO users (name, email, password, role,created_at) VALUES ($1, $2, $3,$4, NOW()) RETURNING id`
	err = r.db.QueryRow(query, user.Name, user.Email, user.Password, user.Role).Scan(&user.ID)
	return err
}

func (r *userRepo) GetUserByID(id string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, password, role, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepo) UpdateUser(user *model.User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, password = $3, role = $4, updated_at = CURRENT_TIMESTAMP
	WHERE id = $5
`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.Role, user.ID)
	return err
}
func (r *userRepo) IsEmailTaken(email string, excludeUserID int) (bool, error) {
	var existingID int
	query := `SELECT id FROM users WHERE email = $1 AND id != $2`
	err := r.db.QueryRow(query, email, excludeUserID).Scan(&existingID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *userRepo) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
func (r *userRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, password, role, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
