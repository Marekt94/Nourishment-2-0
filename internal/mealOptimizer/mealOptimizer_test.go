package mealOptimizer

import (
	"encoding/json"
	meal "nourishment_20/internal/mealDomain"
	"testing"
)

func TestGenSchema(t *testing.T) {
	//given
	expectedSchema := `{
    "type": "object",
    "properties": {
        "products": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "number"
                    },
                    "name": {
                        "type": "string"
                    },
                    "finalweightAfterOptimization": {
                        "type": "number"
                    }
                },
                "additionalProperties": false,
                "required": [
                    "id",
                    "name",
                    "finalweightAfterOptimization"
                ]
            }
        },
        "cumulativeKcal": {
            "type": "number"
        }
    },
    "additionalProperties": false,
    "required": [
        "products",
        "cumulativeKcal"
    ]
}`
	expected := make(map[string]interface{})
	err := json.Unmarshal([]byte(expectedSchema), &expected)
	if err != nil {
		t.Fatalf("Error unmarshaling expected schema: %v", err)
	}
	t.Logf("Expected schema: %v", expected)
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Error marshaling expected schema: %v", err)
	}

	//when
	schema := ProdsInMealResponse{}.GetAIResponseSchema()
	if err != nil {
		t.Fatalf("Error marshaling schema: %v", err)
	}
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Error marshaling returned schema: %v", err)
	}

	//then
	if string(expectedJSON) != string(schemaJSON) {
		t.Errorf("Schema mismatch.\nExpected:\n%s\nGot:\n%s", expectedJSON, schemaJSON)
	}
}

func TestProdToString(t *testing.T) {
	//given
	dest := `- płatki owsiane:
  - id: 13
  - 336 kcal/100g
  - waga wstępna: 60g`

	prod := meal.ProductInMeal{
		Product: meal.Product{Id: 13,
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

	meal := meal.Meal{
		ProductsInMeal: []meal.ProductInMeal{
			{
				Product: meal.Product{
					Id:         13,
					Name:       "płatki owsiane",
					KcalPer100: 336,
					Unit:       "g",
				},
				Weight: 60,
			},
			{
				Product: meal.Product{
					Id:         14,
					Name:       "mleko 2%",
					KcalPer100: 51,
					Unit:       "ml",
				},
				Weight: 300,
			},
			{
				Product: meal.Product{
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
	res := MealToString(&meal)
	//then
	if res != dest {
		t.Errorf("Expected\n%s, got\n%s", dest, res)
	}
}
