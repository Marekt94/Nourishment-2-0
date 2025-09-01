package api

// Resource names - nazwy zasobów
const (
	RESOURCE_LOGIN              = "login"
	RESOURCE_MEALS              = "meals"
	RESOURCE_MEALSINDAY         = "mealsinday"
	RESOURCE_PRODUCTS           = "products"
	RESOURCE_LOOSEPRODUCTSINDAY = "looseproductsinday"
	RESOURCE_CATEGORIES         = "categories"
	RESOURCE_OPTIMIZEMEAL       = "optimizemeal"
)

// API endpoint paths - stałe dla ścieżek endpointów
const (
	// Authentication
	PATH_LOGIN = RESOURCE_LOGIN

	// Meals endpoints
	PATH_MEALS         = "/" + RESOURCE_MEALS
	PATH_MEALS_WITH_ID = "/" + RESOURCE_MEALS + "/:id"

	// Meals in day endpoints
	PATH_MEALSINDAY         = "/" + RESOURCE_MEALSINDAY
	PATH_MEALSINDAY_WITH_ID = "/" + RESOURCE_MEALSINDAY + "/:id"

	// Products endpoints
	PATH_PRODUCTS         = "/" + RESOURCE_PRODUCTS
	PATH_PRODUCTS_WITH_ID = "/" + RESOURCE_PRODUCTS + "/:id"

	// Loose products in day endpoints
	PATH_LOOSEPRODUCTSINDAY         = "/" + RESOURCE_LOOSEPRODUCTSINDAY
	PATH_LOOSEPRODUCTSINDAY_WITH_ID = "/" + RESOURCE_LOOSEPRODUCTSINDAY + "/:id"

	// Categories endpoints
	PATH_CATEGORIES         = "/" + RESOURCE_CATEGORIES
	PATH_CATEGORIES_WITH_ID = "/" + RESOURCE_CATEGORIES + "/:id"

	// Optimization endpoints
	PATH_OPTIMIZEMEAL         = "/" + RESOURCE_OPTIMIZEMEAL
	PATH_OPTIMIZEMEAL_WITH_ID = "/" + RESOURCE_OPTIMIZEMEAL + "/:id"
)
