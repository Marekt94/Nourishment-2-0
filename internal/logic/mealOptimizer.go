package logic

import (
	"fmt"
	utils "nourishment_20/internal"
	"nourishment_20/internal/AIClient"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
	"os"
	"strconv"
	"strings"

	"github.com/invopop/jsonschema"
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

// ToJSONSchema returns the JSON schema for OpenRouterOptimizationResult using github.com/invopop/jsonschema
func ProdsInMealResponseSchema() *jsonschema.Schema {
	schema := jsonschema.Reflect(&ProdsInMealResponse{})
	schema.AdditionalProperties = jsonschema.FalseSchema
	return schema
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
	res, _ = o.AIClient.ExecutePrompt(prompt, ProdsInMealResponseSchema())
	logging.Global.Tracef("AI response:\n%s", res)

	return nil, nil
}
