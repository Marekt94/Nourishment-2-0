package auth

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
)

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"-"` // Hidden from JSON output
}

// UserRepo defines the interface for user repository operations
type UserRepo interface {
	CreateUser(u *User) int64
	DeleteUser(id int) bool
	IsUserExists(username, password string) bool
}

// FirebirdUserRepo implements UserRepo interface using Firebird database
type FirebirdUserRepo struct {
	Database *sql.DB
}

// CreateUser creates a new user and returns the user ID
func (repo *FirebirdUserRepo) CreateUser(u *User) int64 {
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (%s)", USER_TAB, USER_USERNAME, USER_PASSWORD, database.QuestionMarks(2))

	result, err := repo.Database.Exec(sqlStr, u.Username, u.Password)
	if err != nil {
		logging.Global.Warnf("Failed to create user: %v", err)
		return -1
	}

	id, err := result.LastInsertId()
	if err != nil {
		logging.Global.Warnf("Failed to get last insert ID: %v", err)
		return -1
	}

	logging.Global.Infof("User created successfully with ID: %d", id)
	return id
}

// DeleteUser deletes a user by ID and returns true if successful
func (repo *FirebirdUserRepo) DeleteUser(id int) bool {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s = %s", USER_TAB, USER_ID, database.QuestionMarks(1))

	result, err := repo.Database.Exec(sqlStr, id)
	if err != nil {
		logging.Global.Warnf("Failed to delete user: %v", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logging.Global.Warnf("Failed to get rows affected: %v", err)
		return false
	}

	if rowsAffected == 0 {
		logging.Global.Warnf("No user found with ID: %d", id)
		return false
	}

	logging.Global.Infof("User with ID %d deleted successfully", id)
	return true
}

// IsUserExists checks if a user exists with the given username and password
func (repo *FirebirdUserRepo) IsUserExists(username, password string) bool {
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = %s AND %s = %s",
		USER_TAB, USER_USERNAME, database.QuestionMarks(1), USER_PASSWORD, database.QuestionMarks(1))

	var count int
	err := repo.Database.QueryRow(sqlStr, username, password).Scan(&count)
	if err != nil {
		logging.Global.Warnf("Failed to check if user exists: %v", err)
		return false
	}

	exists := count > 0
	logging.Global.Debugf("User exists check for username '%s': %t", username, exists)
	return exists
}
