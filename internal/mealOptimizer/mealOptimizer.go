package mealOptimizer

import (
	"encoding/json"
	"fmt"
	utils "nourishment_20/internal"
	"nourishment_20/internal/AIClient"
	"nourishment_20/internal/logging"
	meal "nourishment_20/internal/mealDomain"
	"os"
	"strconv"
	"strings"
)

type Optimizer struct {
	AIClient AIClient.AIClientIntf
}

type ProdsInMealResponse struct {
	Products       []Product `json:"products"`
	CumulativeKcal float64   `json:"cumulativeKcal"`
}

type Product struct {
	ID                           float64 `json:"id"`
	Name                         string  `json:"name"`
	FinalWeightAfterOptimization float64 `json:"finalweightAfterOptimization"`
}

func (p Product) GetAIResponseSchema() map[string]interface{} {
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

func (p Product) MarshalJSON() ([]byte, error) {
	schema := p.GetAIResponseSchema()
	return json.Marshal(schema)
}

func (p ProdsInMealResponse) GetAIResponseSchema() map[string]interface{} {
	res := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"products": map[string]interface{}{
				"type":  "array",
				"items": Product{}.GetAIResponseSchema(),
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

func findProdInMealDB(productsInMeal []meal.ProductInMeal, index float64) *meal.ProductInMeal {
	for i := range productsInMeal {
		if productsInMeal[i].Product.Id == int(index) {
			return &productsInMeal[i]
		}
	}
	return nil
}
func (p ProdsInMealResponse) UpdateProductsInMeal(productsInMeal []meal.ProductInMeal) {
	for i := range p.Products {
		prodInMealDB := findProdInMealDB(productsInMeal, p.Products[i].ID)
		if prodInMealDB != nil {
			prodInMealDB.Weight = p.Products[i].FinalWeightAfterOptimization
		} else {
			logging.Global.Panicf("Product with ID %v not found in productsInMeal", p.Products[i].ID)
		}
	}
}

func ProdToString(p meal.ProductInMeal) string {
	cProdStringSchema := "- %s:\n  - id: %d\n  - %s kcal/100%s\n  - waga wstępna: %v%s"
	prodString := fmt.Sprintf(cProdStringSchema, p.Product.Name, p.Product.Id, strconv.FormatFloat(p.Product.KcalPer100, 'f', -1, 64),
		p.Product.Unit, p.Weight, p.Product.Unit)
	logging.Global.Debugf(prodString)
	return prodString
}

func MealToString(m *meal.Meal) string {
	prodsInMealStr := []string{}
	for _, prodsInMeal := range m.ProductsInMeal {
		prodsInMealStr = append(prodsInMealStr, ProdToString(prodsInMeal))
	}
	return strings.Join(prodsInMealStr, "\n")
}

func (o *Optimizer) executePromptForOptimizingMeal(m *meal.Meal, k float64) (*string, bool) {
	fileContent, err := utils.ReadFile(AI_OPTIMIZATION_PROMPT)
	if err != nil {
		return nil, false
	}
	logging.Global.Tracef("Read file content:\n%s", fileContent)
	promptScheme := string(fileContent)
	prompt := os.Expand(promptScheme, func(key string) string {
		switch key {
		case MEAL_INGREDIENTS:
			return MealToString(m)
		case KCAL:
			return fmt.Sprintf("%.1f", k)
		}
		return os.Getenv(key)
	})
	logging.Global.Tracef("Optimizing prompt:\n%s", prompt)
	res, succ := o.AIClient.ExecutePrompt(prompt, nil)
	logging.Global.Tracef("AI response:\n%s", res)
	return &res, succ
}

func (o *Optimizer) executePromptForRetrivingOptimizedMealValues(a *string) (*string, bool) {
	fileContent, err := utils.ReadFile(AI_GET_OPTIMIZED_MEAL_PROMPT)
	if err != nil {
		return nil, false
	}
	logging.Global.Tracef("Read file content:\n%s", fileContent)
	promptScheme := string(fileContent)
	prompt := os.Expand(promptScheme, func(key string) string {
		switch key {
		case OPTIMIZATION_ANSWER:
			return *a
		}
		return os.Getenv(key)
	})
	res, succ := o.AIClient.ExecutePrompt(prompt, ProdsInMealResponse{})
	logging.Global.Tracef("AI response:\n%s", res)
	return &res, succ
}

func (o *Optimizer) OptimizeMeal(m *meal.Meal, k float64) (*meal.Meal, error) {
	res, _ := o.executePromptForOptimizingMeal(m, k)
	if res == nil {
		logging.Global.Panicf("Error while executing prompt for optimizing meal. Meal:%v/n", m)
	}

	prodsInMealString, _ := o.executePromptForRetrivingOptimizedMealValues(res)
	if res == nil {
		logging.Global.Panicf("Error while executing prompt to retrive optimized values from AI response. Response: %v\n", res)
	}

	var prodsInMealResponse ProdsInMealResponse
	err := json.Unmarshal([]byte(*prodsInMealString), &prodsInMealResponse)
	if err != nil {
		logging.Global.Panicf("Error unmarshaling AI response: %v", err)
	}
	prodsInMealResponse.UpdateProductsInMeal(m.ProductsInMeal)
	logging.Global.Debugf("Updated products in meal: %v", prodsInMealResponse)
	return m, nil
}
