package database

type MealsRepo interface {
	GetMeal(i int) Meal
	GetMeals() []Meal
	DeleteMeal(i int) bool
	CreateMeal(m *Meal) int64
	UpdateMeal(m *Meal)
}

type Meal struct { // [AI REFACTOR]
	Id             int             `json:"id"`
	Name           string          `json:"name"`
	Recipe         string          `json:"recipe"`
	ProductsInMeal []ProductInMeal `json:"productsInMeal"`
}

type ProductsRepo interface {
	GetProduct(i int) Product
	GetProducts() []Product
	CreateProduct(p *Product) int64
	DeleteProduct(i int) bool
	UpdateProduct(p *Product)
}

type Category struct { // [AI REFACTOR]
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct { // [AI REFACTOR]
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
	Category      Category `json:"category"`
}

type ProductInMeal struct { // [AI REFACTOR]
	Id      int     `json:"id"`
	Product Product `json:"product"`
	Weight  float64 `json:"weight"`
}

// [API GEN] DTO do API (mapowanie z MealInDayDb)
type MealInDay struct {
	Id                    int
	Breakfast             Meal
	SecondBreakfast       Meal
	Lunch                 Meal
	Dinner                Meal
	Supper                Meal
	AfternoonSnack        Meal
	For5Days              bool // mapowane z CHAR(1) '1' lub '0'
	FactorBreakfast       float64
	FactorSecondBreakfast float64
	FactorLunch           float64
	FactorDinner          float64
	FactorSupper          float64
	FactorAfternoonSnack  float64
	Name                  string
}

// [API GEN] Interfejs MealsInDayRepo
type MealsInDayRepo interface {
	CreateMealsInDay(m *MealInDay) int64
	GetMealsInDay(id int) MealInDay
	GetMealsInDays() []MealInDay
	DeleteMealsInDay(id int) bool
	UpdateMealsInDay(m *MealInDay)
}

// LooseProductInDay reprezentuje pojedynczy produkt w dniu (tabela PRODUKTY_LUZNE_W_DNIU)
type LooseProductInDay struct {
	Id      int     `json:"id"`
	DayId   int     `json:"dayId"`   // ID_DNIA
	Product Product `json:"product"` // ID_PRODUKTU jako relacja
	Weight  float64 `json:"weight"`  // ILOSC_W_G
}

// Interfejs do operacji na luźnych produktach w dniu
type LooseProductsInDayRepo interface {
	CreateLooseProductInDay(p *LooseProductInDay) int64
	GetLooseProductInDay(id int) LooseProductInDay
	GetLooseProductsInDay(dayId int) []LooseProductInDay
	DeleteLooseProductInDay(id int) bool
	UpdateLooseProductInDay(p *LooseProductInDay)
}
