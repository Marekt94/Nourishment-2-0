package main

/*
TODO:DODAĆ TESTY W POSTMANIE:
 - update potrawy - czy updatują sie produkty?
DONE: dodać crud, dto, repo dla kategorii
DONE: dodać api dla optymalizacji potraw - kalorycznosc zmienna
DONE: dodać crud, dto, repo dla produktów wolnych w dniu
DONE: testy dla luźnych produktów w dniu
TODO: zwracac w responsie potraw w dniu całkowite makro
TODO: dodać endpoint do wydruku, niech przesyła pdfa (albo w markdown)
DONE: jwt, autoryzacja uwierzytelnianie
TODO: dodać weryfikację uprawnień po danych w bazie danych
DONE: stworzyc gotowego maina, zeby byl wystawialny w prosty sposób
TODO: update postmana
TODO: refactoring globalny
DONE: swagger
*/

import (
	"nourishment_20/internal/modules"
)

// @title           Nourishment 2.0 API
// @version         1.0
// @contact.email   marekt94@gmail.com
// @host            localhost:8080
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Get token from authorization request and place in "Value" field "Bearer {token}"

func main() {
	kernel := modules.NewMealKernel()
	kernel.Init()
	kernel.Run()
}
