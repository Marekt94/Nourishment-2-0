package logic

import (
	"encoding/json"
	"fmt"
	utils "nourishment_20/internal"
	"nourishment_20/internal/AIClient"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
	"os"
	"strconv"
	"strings"
)

type MealOptimizer struct {
	AIClient AIClient.AIClientIntf
}

// OpenRouterOptimizationResult represents the expected JSON structure for OpenRouter response
// swagger:model
type ProdsInMealResponse struct {
	Products       []OpenRouterProduct `json:"products" jsonschema:"required"`
	CumulativeKcal float64             `json:"cumulativeKcal" jsonschema:"required"`
}

type OpenRouterProduct struct {
	ID                           float64 `json:"id" jsonschema:"required"`
	Name                         string  `json:"name" jsonschema:"required"`
	FinalWeightAfterOptimization float64 `json:"finalweightAfterOptimization" jsonschema:"required"`
}

func (p OpenRouterProduct) GetAIResponseSchema() map[string]interface{} {
	res := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id":                           map[string]string{"type": "number"},
			"name":                         map[string]string{"type": "string"},
			"finalweightAfterOptimization": map[string]string{"type": "number"},
		},
		"additionalProperties": false,
		"required":             []string{"id", "name", "finalweightAfterOptimization"},
	}

	return res
}

func (p ProdsInMealResponse) GetAIResponseSchema() map[string]interface{} {
	res := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"products": map[string]interface{}{
				"type":  "array",
				"items": OpenRouterProduct{}.GetAIResponseSchema(),
			},
			"cumulativeKcal": map[string]string{"type": "number"},
		},
		"additionalProperties": false,
		"required":             []string{"products", "cumulativeKcal"},
	}
	return res
}

func (p ProdsInMealResponse) MarshalJSON() ([]byte, error) {
	schema := p.GetAIResponseSchema()
	return json.Marshal(schema)
}

func (p OpenRouterProduct) MarshalJSON() ([]byte, error) {
	schema := p.GetAIResponseSchema()
	return json.Marshal(schema)
}

func ProdToString(p database.ProductInMeal) string {
	cProdStringSchema := "- %s:\n  - id: %d\n  - %s kcal/100%s\n  - waga wstępna: %v%s"
	prodString := fmt.Sprintf(cProdStringSchema, p.Product.Name, p.Product.Id, strconv.FormatFloat(p.Product.KcalPer100, 'f', -1, 64),
		p.Product.Unit, p.Weight, p.Product.Unit)
	logging.Global.Debugf(prodString)
	return prodString
}

func MealToString(m database.Meal) string {
	prodsInMealStr := []string{}
	for _, prodsInMeal := range m.ProductsInMeal {
		prodsInMealStr = append(prodsInMealStr, ProdToString(prodsInMeal))
	}
	return strings.Join(prodsInMealStr, "\n")
}

func (o *MealOptimizer) OptimizeMeal(m database.Meal) (*database.Meal, error) {
	fileContent, err := utils.ReadFile("ai_optimization_prompt.md")
	if err != nil {
		return nil, err
	}
	logging.Global.Tracef("Read file content:\n%s", fileContent)
	promptScheme := string(fileContent)
	prompt := os.Expand(promptScheme, func(key string) string {
		switch key {
		case "MEAL_INGREDIENTS":
			return MealToString(m)
		}
		return os.Getenv(key)
	})
	logging.Global.Tracef("Optimizing prompt:\n%s", prompt)
	res, _ := o.AIClient.ExecutePrompt(prompt, nil)
	logging.Global.Tracef("AI response:\n%s", res)

	fileContent, err = utils.ReadFile("ai_get_optimized_ingredients_prompt.md")
	if err != nil {
		return nil, err
	}
	logging.Global.Tracef("Read file content:\n%s", fileContent)
	promptScheme = string(fileContent)
	prompt = os.Expand(promptScheme, func(key string) string {
		switch key {
		case "OPTIMIZATION_ANSWER":
			return res
		}
		return os.Getenv(key)
	})
	res, _ = o.AIClient.ExecutePrompt(prompt, ProdsInMealResponse{})
	logging.Global.Tracef("AI response:\n%s", res)

	return nil, nil
}
