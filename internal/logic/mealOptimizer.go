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
)

type MealOptimizer struct {
	AIClient AIClient.AIClientIntf
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
	logging.Global.Debugf("Read file content:\n%s", fileContent)
	optimizingPromptScheme := string(fileContent)
	optimizingPrompt := os.Expand(optimizingPromptScheme, func(key string) string {
		switch key {
		case "MEAL_INGREDIENTS":
			return MealToString(m)
		}
		return os.Getenv(key)
	})
	logging.Global.Debugf("Optimizing prompt:\n%s", optimizingPrompt)
	res, _ := o.AIClient.ExecutePrompt(optimizingPrompt)
	logging.Global.Debugf("AI response:\n%s", res)
	return nil, nil
}
