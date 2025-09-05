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
TODO: swagger
*/

import (
	"nourishment_20/internal/modules"
)

func StartMealServer() {
	kernel := modules.NewMealKernel()
	kernel.Init()
	kernel.Run()
}

func main() {
	StartMealServer() // [AI REFACTOR] uruchom serwer
}
