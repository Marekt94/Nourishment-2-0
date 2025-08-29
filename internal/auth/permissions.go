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
}

type PermissionsRepo struct {
	db *sql.DB
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
	sqlStr := fmt.Sprintf("SELECT p.%s, p.%s FROM %s up JOIN %s p ON up.%s = p.%s WHERE up.%s = %s",
		PERMISSIONS_RESOURCE, PERMISSIONS_RIGHT,
		USER_PERMISSIONS_TAB, PERMISSIONS_TAB,
		USER_PERMISSIONS_PERMISSION_ID, PERMISSIONS_ID,
		USER_PERMISSIONS_USER, database.QuestionMarks(1))

	rows, err := p.db.Query(sqlStr, user)
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
		err := p.db.QueryRow(checkSql, resource, right).Scan(&existingID)
		if err != nil && err != sql.ErrNoRows {
			logging.Global.Warnf("Failed to check existing permission: %v", err)
			continue
		}

		// Jeśli uprawnienie nie istnieje, dodaj je
		if err == sql.ErrNoRows {
			insertSql := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (%s) RETURNING %s",
				PERMISSIONS_TAB, PERMISSIONS_RESOURCE, PERMISSIONS_RIGHT, database.QuestionMarks(2), PERMISSIONS_ID)

			var newID int64
			err := p.db.QueryRow(insertSql, resource, right).Scan(&newID)
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
