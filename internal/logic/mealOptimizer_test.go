package logic

import (
	"nourishment_20/internal/database"
	"testing"
)

func TestProdToString(t *testing.T) {
	//given
	dest := `- płatki owsiane:
  - id: 13
  - 336 kcal/100g
  - waga wstępna: 60g`

	prod := database.ProductInMeal{
		Product: database.Product{Id: 13,
			Name:       "płatki owsiane",
			KcalPer100: 336,
			Unit:       "g",
		},
		Weight: 60,
	}
	//when
	res := ProdToString(prod)
	//then
	if res != dest {
		t.Errorf("Expected\n%s, got\n%s", dest, res)
	}
}

func TestMealToString(t *testing.T) {
	//given
	dest := `- płatki owsiane:
  - id: 13
  - 336 kcal/100g
  - waga wstępna: 60g
- mleko 2%:
  - id: 14
  - 51 kcal/100ml
  - waga wstępna: 300ml
- banan:
  - id: 12
  - 97 kcal/100g
  - waga wstępna: 120g`

	meal := database.Meal{
		ProductsInMeal: []database.ProductInMeal{
			{
				Product: database.Product{
					Id:         13,
					Name:       "płatki owsiane",
					KcalPer100: 336,
					Unit:       "g",
				},
				Weight: 60,
			},
			{
				Product: database.Product{
					Id:         14,
					Name:       "mleko 2%",
					KcalPer100: 51,
					Unit:       "ml",
				},
				Weight: 300,
			},
			{
				Product: database.Product{
					Id:         12,
					Name:       "banan",
					KcalPer100: 97,
					Unit:       "g",
				},
				Weight: 120,
			},
		},
	}

	//when
	res := MealToString(meal)
	//then
	if res != dest {
		t.Errorf("Expected\n%s, got\n%s", dest, res)
	}
}
