package models

import (
	"database/sql"
	"fmt"
	"time"

	"gosql_interface/internal/utils"
)

// User represents a user in the system
type User struct {
	ID        int64      `json:"id"`
	Name      *string    `json:"name"`
	Login     string     `json:"login"`
	Password  string     `json:"password,omitempty"` // Input only, cleared after operations
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
}

// GetAllUsers retrieves all users with optional search
func GetAllUsers(db *sql.DB, search string) ([]User, error) {
	query := `
		SELECT id, name, login, password, role, created_at
		FROM users
		WHERE ($1 = '' OR name ILIKE $2 OR login ILIKE $2)
		ORDER BY id DESC
	`

	searchPattern := "%" + search + "%"
	rows, err := db.Query(query, search, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Login, &u.Password, &u.Role, &u.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	query := `SELECT id, name, login, password, role, created_at FROM users WHERE id = $1`

	var u User
	err := db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Login, &u.Password, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

// GetUserByLogin retrieves a user by login (for authentication)
func GetUserByLogin(db *sql.DB, login string) (*User, error) {
	query := `SELECT id, name, login, password, role, created_at FROM users WHERE login = $1`

	var u User
	err := db.QueryRow(query, login).Scan(&u.ID, &u.Name, &u.Login, &u.Password, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

// CreateUser creates a new user with hashed password
func CreateUser(db *sql.DB, user *User) error {
	// Hash password before storing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (name, login, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err = db.QueryRow(query, user.Name, user.Login, hashedPassword, user.Role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user
func UpdateUser(db *sql.DB, user *User) error {
	// If password is provided, hash it
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		query := `
			UPDATE users
			SET name = $1, login = $2, password = $3, role = $4
			WHERE id = $5
		`
		_, err = db.Exec(query, user.Name, user.Login, hashedPassword, user.Role, user.ID)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		// Update without changing password
		query := `
			UPDATE users
			SET name = $1, login = $2, role = $3
			WHERE id = $4
		`
		_, err := db.Exec(query, user.Name, user.Login, user.Role, user.ID)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

// DeleteUser deletes a user by ID
func DeleteUser(db *sql.DB, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
