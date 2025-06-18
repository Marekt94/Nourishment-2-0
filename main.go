package main

import (
	"fmt"
	database "nourishment_20/data"
)

func main() {
	var conf database.DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:\Users\marek\Documents\nourishment_backup_db\NOURISHMENT.FDB`
    
	fDbEngine := database.FBDBEngine{BaseEngineIntf: &database.BaseEngine{}}

	engine := fDbEngine.Connect(&conf)

	var mealsRepo database.MealsRepo = &database.FirebirdRepoAccess{DbEngine: engine}
	fmt.Println(mealsRepo.GetMeal(88))
}