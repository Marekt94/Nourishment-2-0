package meal

import (
	"testing"

	"nourishment_20/internal/database"
)

const TEST_CAT_NAME = "test_category"

func initCategoryRepo() CategoriesRepo {
	var conf database.DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:\\Users\\marek\\Documents\\nourishment_backup_db\\NOURISHMENT.FDB`

	fDbEngine := FBDBEngine{BaseEngineIntf: &database.BaseEngine{}}
	engine := fDbEngine.Connect(&conf)
	return &FirebirdRepoAccess{Database: engine}
}

func createCategory() Category {
	repo := initCategoryRepo()
	c := Category{Name: TEST_CAT_NAME}
	id := repo.CreateCategory(&c)
	c.Id = int(id)
	return c
}

func TestCreateCategory(t *testing.T) {
	c := createCategory()
	if c.Id <= 0 {
		t.Error("no category created")
	}
	if c.Name != TEST_CAT_NAME {
		t.Error("wrong name in created category")
	}
}

func TestGetCategory(t *testing.T) {
	repo := initCategoryRepo()
	c := createCategory()
	got := repo.GetCategory(c.Id)
	if got.Id != c.Id {
		t.Error("fail to get category")
	}
}

func TestGetCategories(t *testing.T) {
	repo := initCategoryRepo()
	cats := repo.GetCategories()
	if len(cats) < 1 {
		t.Error("no categories retrieved")
	}
}

func TestUpdateCategory(t *testing.T) {
	repo := initCategoryRepo()
	c := createCategory()
	newName := "updated_cat"
	c.Name = newName
	repo.UpdateCategory(&c)
	got := repo.GetCategory(c.Id)
	if got.Name != newName {
		t.Error("category not updated")
	}
}

func TestDeleteCategory(t *testing.T) {
	repo := initCategoryRepo()
	c := createCategory()
	ok := repo.DeleteCategory(c.Id)
	if !ok {
		t.Errorf("category with id %d not deleted", c.Id)
	}
}
