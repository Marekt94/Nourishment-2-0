package auth

import (
	log "nourishment_20/internal/logging"
	h "nourishment_20/internal/testHelper"
	"testing"
)

const TEST_USER_NAME = "testuser"
const TEST_USER_PASSWORD = "testpassword"

var REPO = createRepo()

func createRepo() UserRepo {
	c := h.InitTestUnit(`..\..\.env`)
	return &FirebirdUserRepo{Database: c}
}

func setupTestLogin(t *testing.T) int64 {
	log.Global.Infof("Setting up TestLogin...")
	user := User{
		Username: TEST_USER_NAME,
		Password: TEST_USER_PASSWORD,
	}
	return REPO.CreateUser(&user)
}

func teardownTestLogin(t *testing.T, i int) {
	log.Global.Infof("Tearing down TestLogin...")
	REPO.DeleteUser(i)
}

func TestLogin(t *testing.T) {
	i := setupTestLogin(t)
	if i <= 0 {
		t.Error("User not created")
	}
	userId := REPO.IsUserExists(TEST_USER_NAME, TEST_USER_PASSWORD)
	defer teardownTestLogin(t, int(i))
	if userId <= 0 {
		t.Error("User should exist after creation")
	}
	if userId != i {
		t.Errorf("Expected user ID %d, got %d", i, userId)
	}
	nonExistentUserId := REPO.IsUserExists("RAND", "RAND")
	if nonExistentUserId != -1 {
		t.Error("User should not exist with random credentials, expected -1")
	}
}
