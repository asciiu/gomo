package sql_test

import (
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

// func setupService() *sql.DB {
// 	dbUrl := "postgres://fomo:fomornd@127.0.0.1:6432/fomo_dev?&sslmode=disable"
// 	db, error := db.NewDB(dbUrl)

// 	//user := user.NewUser("first", "last", "test@email", "hash")
// 	//_, error := userRepo.InsertUser(db, user)
// 	checkErr(error)

// 	return db
// }

// func TestFindBalances(t *testing.T) {
// 	db := setupService()
// 	defer db.Close()

// 	ding := bp.GetUserBalancesRequest{
// 		UserId:   "bf68c122-37b2-4e63-9605-bcae7b82b04f",
// 		ApiKeyId: "71672e76-8515-4785-b931-bdbf6157b198",
// 	}

// 	balances, error := balRepo.FindBalancesByUserId(db, &ding)
// 	if error != nil {
// 		fmt.Println(error)
// 	}
// 	fmt.Println(balances)
// }
