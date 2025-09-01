package auth

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
)

type PermissionDB struct {
	Resource string
	Right    string
}

type PermissionPerUserDB struct {
	PermissionDB
	User string
}

type Permission struct {
	Resource string
	Rights   []string
}

type PermissionsIntf interface {
	GetPermissions(user string) []Permission
	RegisterPermissions(resource string, rights []string)
	RegisterUserPermission(username string, resource string, right string) error
}

type PermissionsRepo struct {
	Db *sql.DB
}

type PermissionController struct {
	Repo PermissionsIntf
}

func (c *PermissionController) convertToPermissions(s *map[string][]string) []Permission {
	// Konwersja na slice Permission

	var permissions []Permission
	for resource, rights := range *s {
		permissions = append(permissions, Permission{
			Resource: resource,
			Rights:   rights,
		})
	}

	return permissions
}

// TODO pobieraj prawa dla danego usera z bazy danych
func (p *PermissionsRepo) GetPermissions(user string) []Permission {
	// Zapytanie z JOIN do tabeli użytkowników, bo USER_PERMISSIONS_USER to teraz FK do UZYTKOWNICY.ID
	sqlStr := fmt.Sprintf(`SELECT p.%s, p.%s 
		FROM %s up 
		JOIN %s p ON up.%s = p.%s 
		JOIN %s u ON up.%s = u.%s 
		WHERE u.%s = %s`,
		PERMISSIONS_RESOURCE, PERMISSIONS_RIGHT,
		USER_PERMISSIONS_TAB, PERMISSIONS_TAB, USER_PERMISSIONS_PERMISSION_ID, PERMISSIONS_ID,
		USER_TAB, USER_PERMISSIONS_USER, USER_ID,
		USER_USERNAME, database.QuestionMarks(1))

	rows, err := p.Db.Query(sqlStr, user)
	if err != nil {
		logging.Global.Warnf("Failed to query user permissions: %v", err)
		return nil
	}
	defer rows.Close()

	// Kolekcja PermissionDB z bazy danych
	var permissionsDB []PermissionDB
	for rows.Next() {
		var permDB PermissionDB
		if err := rows.Scan(&permDB.Resource, &permDB.Right); err != nil {
			logging.Global.Warnf("Failed to scan permission row: %v", err)
			continue
		}
		permissionsDB = append(permissionsDB, permDB)
	}

	if err := rows.Err(); err != nil {
		logging.Global.Warnf("Error iterating over permission rows: %v", err)
		return nil
	}

	// Grupowanie PermissionDB po Resource
	permissionsMap := make(map[string][]string)
	for _, permDB := range permissionsDB {
		permissionsMap[permDB.Resource] = append(permissionsMap[permDB.Resource], permDB.Right)
	}

	// Konwersja na kolekcję Permission
	var permissions []Permission
	for resource, rights := range permissionsMap {
		permissions = append(permissions, Permission{
			Resource: resource,
			Rights:   rights,
		})
	}

	logging.Global.Debugf("Retrieved %d permissions for user '%s'", len(permissions), user)
	return permissions
}

func (p *PermissionsRepo) RegisterPermissions(resource string, rights []string) {
	for _, right := range rights {
		// Sprawdź czy uprawnienie już istnieje
		checkSql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s AND %s = %s",
			PERMISSIONS_ID, PERMISSIONS_TAB, PERMISSIONS_RESOURCE, database.QuestionMarks(1),
			PERMISSIONS_RIGHT, database.QuestionMarks(1))

		var existingID int64
		err := p.Db.QueryRow(checkSql, resource, right).Scan(&existingID)
		if err != nil && err != sql.ErrNoRows {
			logging.Global.Warnf("Failed to check existing permission: %v", err)
			continue
		}

		// Jeśli uprawnienie nie istnieje, dodaj je
		if err == sql.ErrNoRows {
			insertSql := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (%s) RETURNING %s",
				PERMISSIONS_TAB, PERMISSIONS_RESOURCE, PERMISSIONS_RIGHT, database.QuestionMarks(2), PERMISSIONS_ID)

			var newID int64
			err := p.Db.QueryRow(insertSql, resource, right).Scan(&newID)
			if err != nil {
				logging.Global.Warnf("Failed to insert permission '%s:%s': %v", resource, right, err)
			} else {
				logging.Global.Infof("Registered permission '%s:%s' with ID: %d", resource, right, newID)
			}
		} else {
			logging.Global.Debugf("Permission '%s:%s' already exists with ID: %d", resource, right, existingID)
		}
	}
}

// RegisterUserPermission rejestruje uprawnienie dla danego użytkownika
// Parametry: nazwa użytkownika, zasób, uprawnienie (np. "admin", "meals", "read")
func (p *PermissionsRepo) RegisterUserPermission(username string, resource string, right string) error {
	// 1. Znajdź ID użytkownika po nazwie
	userIDSql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s",
		USER_ID, USER_TAB, USER_USERNAME, database.QuestionMarks(1))

	var userID int
	err := p.Db.QueryRow(userIDSql, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user '%s' not found", username)
		}
		logging.Global.Warnf("Failed to find user '%s': %v", username, err)
		return err
	}

	// 2. Znajdź uprawnienie (zasób + prawo) - nie twórz jeśli nie istnieje
	var permissionID int64
	checkPermSql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s AND %s = %s",
		PERMISSIONS_ID, PERMISSIONS_TAB, PERMISSIONS_RESOURCE, database.QuestionMarks(1),
		PERMISSIONS_RIGHT, database.QuestionMarks(1))

	err = p.Db.QueryRow(checkPermSql, resource, right).Scan(&permissionID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Uprawnienie nie istnieje - zwróć błąd
			return fmt.Errorf("permission '%s:%s' does not exist", resource, right)
		}
		logging.Global.Warnf("Failed to check permission '%s:%s': %v", resource, right, err)
		return err
	}

	// 3. Sprawdź czy użytkownik już ma to uprawnienie
	checkUserPermSql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = %s AND %s = %s",
		USER_PERMISSIONS_TAB, USER_PERMISSIONS_USER, database.QuestionMarks(1),
		USER_PERMISSIONS_PERMISSION_ID, database.QuestionMarks(1))

	var count int
	err = p.Db.QueryRow(checkUserPermSql, userID, permissionID).Scan(&count)
	if err != nil {
		logging.Global.Warnf("Failed to check existing user permission: %v", err)
		return err
	}

	// 4. Jeśli użytkownik nie ma tego uprawnienia, dodaj je
	if count == 0 {
		insertUserPermSql := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (%s)",
			USER_PERMISSIONS_TAB, USER_PERMISSIONS_USER, USER_PERMISSIONS_PERMISSION_ID, database.QuestionMarks(2))

		_, err := p.Db.Exec(insertUserPermSql, userID, permissionID)
		if err != nil {
			logging.Global.Warnf("Failed to assign permission '%s:%s' to user '%s': %v", resource, right, username, err)
			return err
		}
		logging.Global.Infof("Successfully assigned permission '%s:%s' to user '%s'", resource, right, username)
	} else {
		logging.Global.Debugf("User '%s' already has permission '%s:%s'", username, resource, right)
	}

	return nil
}
