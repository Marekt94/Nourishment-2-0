package main

import (
	"fmt"
	"log"
	"os"

	"nourishment_20/internal/database"
	"nourishment_20/internal/mealDomain"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DbEngine := database.FBDBEngine{BaseEngineIntf: &database.BaseEngine{}}
	conf := database.DBConf{
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Address:    os.Getenv("DB_ADDRESS"),
		PathOrName: os.Getenv("DB_NAME"),
	}
	db := DbEngine.Connect(&conf)

	repo := meal.FirebirdRepoAccess{Database: db}
	
	// Get products
	products := repo.GetProducts()
	fmt.Printf("Number of products in DB: %d\n", len(products))
	
	fmt.Println("First few products:")
	for i := 0; i < len(products) && i < 5; i++ {
		fmt.Printf("- ID: %d, Name: %s\n", products[i].Id, products[i].Name)
	}
	
	// Check specifically for product ID 2 which failed in the log
	p2 := repo.GetProduct(2)
	fmt.Printf("\nProduct ID 2: ID=%d, Name=%s\n", p2.Id, p2.Name)
}
