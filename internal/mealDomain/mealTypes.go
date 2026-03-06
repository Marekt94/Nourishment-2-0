package meal

type MealsRepoIntf interface {
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

type ProductsRepoIntf interface {
	GetProduct(i int) Product
	GetProducts() []Product
	CreateProduct(p *Product) int64
	DeleteProduct(i int) bool
	UpdateProduct(p *Product)
}

// CategoriesRepoIntf defines CRUD for product categories
type CategoriesRepoIntf interface {
	GetCategory(i int) Category
	GetCategories() []Category
	CreateCategory(c *Category) int64
	DeleteCategory(i int) bool
	UpdateCategory(c *Category)
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
	Id                    int     `json:"id"`
	Breakfast             Meal    `json:"breakfast"`
	SecondBreakfast       Meal    `json:"secondBreakfast"`
	Lunch                 Meal    `json:"lunch"`
	Dinner                Meal    `json:"dinner"`
	Supper                Meal    `json:"supper"`
	AfternoonSnack        Meal    `json:"afternoonSnack"`
	For5Days              bool    `json:"for5Days"` // mapowane z CHAR(1) '1' lub '0'
	FactorBreakfast       float64 `json:"factorBreakfast"`
	FactorSecondBreakfast float64 `json:"factorSecondBreakfast"`
	FactorLunch           float64 `json:"factorLunch"`
	FactorDinner          float64 `json:"factorDinner"`
	FactorSupper          float64 `json:"factorSupper"`
	FactorAfternoonSnack  float64             `json:"factorAfternoonSnack"`
	Name                  string              `json:"name"`
	LooseProducts         []LooseProductInDay `json:"looseProducts"`
}

// [API GEN] Interfejs MealsInDayRepoIntf
type MealsInDayRepoIntf interface {
	CreateMealsInDay(m *MealInDay) int64
	GetMealsInDay(id int) MealInDay
	GetMealsInDays() []MealInDay
	DeleteMealsInDay(id int) bool
	UpdateMealsInDay(m *MealInDay) bool
}

// LooseProductInDay reprezentuje pojedynczy produkt w dniu (tabela PRODUKTY_LUZNE_W_DNIU)
type LooseProductInDay struct {
	Id      int     `json:"id"`
	DayId   int     `json:"dayId"`   // ID_DNIA
	Product Product `json:"product"` // ID_PRODUKTU jako relacja
	Weight  float64 `json:"weight"`  // ILOSC_W_G
}

// Interfejs do operacji na luźnych produktach w dniu
type LooseProductsInDayRepoIntf interface {
	CreateLooseProductInDay(p *LooseProductInDay) int64
	GetLooseProductInDay(id int) LooseProductInDay
	GetLooseProductsInDay(dayId int) []LooseProductInDay
	DeleteLooseProductInDay(id int) bool
	UpdateLooseProductInDay(p *LooseProductInDay)
}

func NewMealInDay() *MealInDay {
	emptyMeal := Meal{Id: EMPTY_ID}

	return &MealInDay{
		Id:                    EMPTY_ID,
		Name:                  "",
		Breakfast:             emptyMeal,
		SecondBreakfast:       emptyMeal,
		Lunch:                 emptyMeal,
		Dinner:                emptyMeal,
		Supper:                emptyMeal,
		AfternoonSnack:        emptyMeal,
		FactorBreakfast:       1.0,
		FactorSecondBreakfast: 1.0,
		FactorLunch:           1.0,
		FactorDinner:          1.0,
		FactorSupper:          1.0,
		FactorAfternoonSnack:  1.0,
		For5Days:              false,
		LooseProducts:         []LooseProductInDay{},
	}
}
