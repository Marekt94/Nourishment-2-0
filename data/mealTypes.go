package database

type MealsRepo interface {
	GetMeal(i int) meal
	GetMeals() []meal
	DeleteMeal(i int) bool
	CreateMeal(m *meal) int64
	UpdateMeal(m *meal)
}

type meal struct {
	Id             int             `json:"id"`
	Name           string          `json:"name"`
	Recipe         string          `json:"recipe"`
	ProductsInMeal []productInMeal `json:"productsInMeal"`
}

type ProductsRepo interface {
	GetProduct(i int) product
	GetProducts() []product
	CreateProduct(p *product) int64
	DeleteProduct(i int) bool
	UpdateProduct(p *product)
}

type category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type product struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	KcalPer100    float64  `json:"kcalPer100"`
	UnitWeight    float64  `json:"weight"`
	Proteins      float64  `json:"proteins"`
	Fat           float64  `json:"fat"`
	Sugar         float64  `json:"sugar"`
	Carbohydrates float64  `json:"carbohydrates"`
	SugarAndCarbo float64  `json:"sugarAndCarb"`
	Fiber         float64  `json:"fiber"`
	Salt          float64  `json:"salt"`
	Unit          string   `json:"unit"`
	Category      category `json:"category"`
}

type productInMeal struct {
	Id      int     `json:"id"`
	Product product `json:"product"`
	Weight  float64 `json:"weight"`
}
