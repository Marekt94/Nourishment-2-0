package meal

import (
	"time"
)

type MealPlanForShoppingList struct {
	MealInDayId int `json:"mealInDayId"`
	Days        int `json:"days"`
}

type LooseProductForShoppingList struct {
	ProductId   int     `json:"productId"`
	ProductName string  `json:"productName"`
	Weight      float64 `json:"weight"`
}

type GenerateShoppingListRequest struct {
	Name          string                        `json:"name"`
	MealPlans     []MealPlanForShoppingList     `json:"mealPlans"`
	LooseProducts []LooseProductForShoppingList `json:"looseProducts"`
}

// GenerateShoppingList creates a shopping list based on selected meal plans and loose products
func GenerateShoppingList(
	req *GenerateShoppingListRequest,
	mealsInDayRepo MealsInDayRepoIntf,
	shoppingListRepo ShoppingListRepoIntf,
	productsRepo ProductsRepoIntf,
) (int64, error) {

	// A map to temporarily aggregate required product quantities by product ID
	// Key: product ID
	aggregatedProducts := make(map[int]float64)

	// A helper func to add weight to product
	addWeight := func(id int, weight float64) {
		aggregatedProducts[id] += weight
	}

	for _, planReq := range req.MealPlans {
		// Fetch MealInDay
		mid := mealsInDayRepo.GetMealsInDay(planReq.MealInDayId)
		if mid.Id == 0 {
			// Skip if not found
			continue
		}

		days := float64(planReq.Days)

		// Aggregate products from each meal type
		aggregateMealProducts := func(m Meal, factor float64) {
			if factor == 0 {
				factor = 1.0
			}
			for _, pim := range m.ProductsInMeal {
				addWeight(pim.Product.Id, pim.Weight*factor*days)
			}
		}

		aggregateMealProducts(mid.Breakfast, mid.FactorBreakfast)
		aggregateMealProducts(mid.SecondBreakfast, mid.FactorSecondBreakfast)
		aggregateMealProducts(mid.Lunch, mid.FactorLunch)
		aggregateMealProducts(mid.Dinner, mid.FactorDinner)
		aggregateMealProducts(mid.Supper, mid.FactorSupper)
		aggregateMealProducts(mid.AfternoonSnack, mid.FactorAfternoonSnack)

		// Aggregate loose products assigned to this MealInDay
		for _, lpid := range mid.LooseProducts {
			addWeight(lpid.Product.Id, lpid.Weight*days)
		}
	}

	// Add manual loose products from the request
	for _, lpReq := range req.LooseProducts {
		addWeight(lpReq.ProductId, lpReq.Weight)
	}

	// Now build the ShoppingList object
	sl := ShoppingList{
		Name:      req.Name,
		CreatedAt: time.Now(),
		EditDate:  time.Now(),
		Products:  []ProductInShoppingList{},
	}

	// For each aggregated product, we need some details, so we fetch it
	for prodId, totalWeight := range aggregatedProducts {
		if totalWeight <= 0 {
			continue
		}

		// It's helpful to populate standard info
		prod := productsRepo.GetProduct(prodId)

		// Ensure the product actually exists in the database before adding to list
		if prod.Id <= 0 {
			continue
		}

		psl := ProductInShoppingList{
			ProductId:    prod.Id,
			ProductName:  prod.Name, // db preferred
			// CategoryName and ProductUnit are mainly mapped in frontend or specific SELECTs, 
			// but we save them here just in case they are used immediately
			CategoryName: prod.Category.Name,
			ProductUnit:  prod.Unit,
			Weight:       totalWeight,
			Bought:       false,
			EditDate:     time.Now(),
		}
		if prod.UnitWeight > 0 {
			psl.Quantity = totalWeight / prod.UnitWeight
		} else {
			psl.Quantity = 0
		}
		sl.Products = append(sl.Products, psl)
	}

	// Save to repo
	id := shoppingListRepo.CreateShoppingList(&sl)

	return id, nil
}
