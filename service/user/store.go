package user

import (
	"database/sql"
	"fmt"

	"github.com/davidado/go-api-reference/types"
)

// Store : User store
type Store struct {
	db *sql.DB
}

// NewStore : Create a new user store
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetUserByEmail : Get user by email
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ? LIMIT 1", email)
	if err != nil {
		return nil, err
	}

	u := &types.User{}
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("user not found")
			}
			return nil, err
		}
	}

	return u, nil
}

// GetUserByID : Get user by ID
func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}

	u := &types.User{}
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("user not found")
			}
			return nil, err
		}
	}

	return u, nil
}

// CreateUser : Create a new user
func (s *Store) CreateUser(u types.User) error {
	_, err := s.db.Exec("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)", u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	u := &types.User{}
	err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
