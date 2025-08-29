package auth

import (
	"reflect"
	"testing"
)

//DONE dorobić permisions repo mock

const TEST_SCOPE_KEY_1 = "meal"
const TEST_SCOPE_VAL_1 = "read"
const TEST_SCOPE_VAL_2 = "write"
const TEST_SCOPE_KEY_2 = "prods"

type PermissionsMock struct {
}

func (p *PermissionsMock) GetPermissions(u string) []Permission {
	return []Permission{
		{Resource: TEST_SCOPE_KEY_1, Rights: []string{TEST_SCOPE_VAL_1, TEST_SCOPE_VAL_2}},
		{Resource: TEST_SCOPE_KEY_2, Rights: []string{TEST_SCOPE_VAL_2}},
	}
}

func (p *PermissionsMock) RegisterPermissions(resource string, rights []string) {

}

func initAuthorizationTestUnit() {
	i := REPO.CreateUser(&User{Username: TEST_USER_NAME, Password: TEST_USER_PASSWORD})
	defer REPO.DeleteUser(int(i))
}

func TestGetJWT(t *testing.T) {
	initAuthorizationTestUnit()
	permissionsRepo := PermissionsMock{}
	jwtGenerator := JWTGenerator{&permissionsRepo}
	permController := PermissionController{&permissionsRepo}
	REPO.IsUserExists(TEST_USER_NAME, TEST_USER_PASSWORD)
	jwt, _ := jwtGenerator.GetJWT("test")
	if jwt == nil {
		t.Fatalf("No JWT")
	}
	tokenString, err := jwtGenerator.JWTToString(jwt)
	if err != nil {
		t.Fatalf("failed to convert JWT to string: %v", err)
	}
	if tokenString == nil {
		t.Fatalf("empty token string")
	}
	actualScopeMap := jwtGenerator.GetScope(*tokenString)
	actualScope := permController.convertToPermissions(&actualScopeMap)
	expectedScope := permissionsRepo.GetPermissions(TEST_USER_NAME)

	if !reflect.DeepEqual(actualScope, expectedScope) {
		t.Fatalf("Expected scope: %v, actual scope: %v", expectedScope, actualScope)
	}

}

//TODO oprogramowac metody do tworzenia string-a z uprawnieniami oraz parsowania uprawnien na mape - dla kilku
