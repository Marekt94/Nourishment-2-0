package database

import "testing"

const TEST_PROD_NAME = "test_name"

func initProductRepo() ProductsRepo {
	var conf DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:\Users\marek\Documents\nourishment_backup_db\NOURISHMENT.FDB`

	fDbEngine := FBDBEngine{BaseEngineIntf: &BaseEngine{}}

	engine := fDbEngine.Connect(&conf)

	return &FirebirdRepoAccess{DbEngine: engine}
}

func createProduct() Product { // [API GEN]
	repo := initProductRepo()
	prod := Product{Name: TEST_PROD_NAME, KcalPer100: 7, UnitWeight: 15, Proteins: 7, Fat: 7, Sugar: 7, Carbohydrates: 7, Fiber: 7, Salt: 1, Unit: `kg`, Category: Category{Id: 1, Name: `no category`}}
	i := int(repo.CreateProduct(&prod))
	prod.Id = i
	return prod
}
func TestCreateProduct(t *testing.T){
	prod := createProduct(); 
	if prod.Id <= 0{
		t.Error(`no product created`)
	}
	if prod.Name != TEST_PROD_NAME{
		t.Error(`wrong name in created product`)
	}
}
func TestGetProduct(t *testing.T) {
	repo := initProductRepo();
	prod := repo.GetProduct(1)
	if (prod.Id != 1) {
		t.Error(`fail to get product`)
	}
}

func TestGetProduct_WhenNoProductFound(t *testing.T) {
	repo := initProductRepo();
	prod := repo.GetProduct(100000)
	if ((prod.Id != 0) || (prod.Name != ``)) {
		t.Error(`fail to get product`)
	}
}

func TestGetProducts(t * testing.T){
	repo := initProductRepo()
	prods := repo.GetProducts()
	t.Logf(`products count: %d`, len(prods))
	if len(prods) < 2 {
		t.Error(`no products retrieved`)
	}
}

func TestDeleteProduct(t *testing.T){
	repo := initProductRepo();
	prod := createProduct();
	res := repo.DeleteProduct(prod.Id)
	if !res {
		t.Errorf(`product with id %d not deleted`, prod.Id)
	}
}

func TestUpdateProduct(t* testing.T){
	newName := `after update`
	newFiber := 1994
	repo := initProductRepo();
	prod := createProduct()
	prod.Name = newName
	prod.Fiber = float64(newFiber)
	repo.UpdateProduct(&prod)
	res := repo.GetProduct(prod.Id)
	if (res.Name != newName) || (res.Fiber != float64(newFiber)){
		t.Error(`product not updated`)
	}
}