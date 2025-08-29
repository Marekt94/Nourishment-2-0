package auth

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
)

// User represents a user in the system
type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Hidden from JSON output
}

// UserRepo defines the interface for user repository operations
type UserRepo interface {
	CreateUser(u *User) int64
	DeleteUser(id int) bool
	IsUserExists(username, password string) int64
}

// FirebirdUserRepo implements UserRepo interface using Firebird database
type FirebirdUserRepo struct {
	Database *sql.DB
}

// CreateUser creates a new user and returns the user ID
func (repo *FirebirdUserRepo) CreateUser(u *User) int64 {
	// Firebird 2.5 doesn't support LastInsertId() properly, so we use RETURNING clause
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (%s) RETURNING %s",
		USER_TAB, USER_USERNAME, USER_PASSWORD, database.QuestionMarks(2), USER_ID)

	var id int64
	err := repo.Database.QueryRow(sqlStr, u.Username, u.Password).Scan(&id)
	if err != nil {
		logging.Global.Warnf("Failed to create user: %v", err)
		return -1
	}

	// Set the ID in the User struct
	u.Id = id

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
// Returns user ID if exists, -1 if not found
func (repo *FirebirdUserRepo) IsUserExists(username, password string) int64 {
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s AND %s = %s",
		USER_ID, USER_TAB, USER_USERNAME, database.QuestionMarks(1), USER_PASSWORD, database.QuestionMarks(1))

	var id int64
	err := repo.Database.QueryRow(sqlStr, username, password).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			logging.Global.Debugf("User not found for username '%s'", username)
			return -1
		}
		logging.Global.Warnf("Failed to check if user exists: %v", err)
		return -1
	}

	logging.Global.Debugf("User found for username '%s' with ID: %d", username, id)
	return id
}
