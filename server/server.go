package main

import "fmt"

func main() {
	// myRouter := mux.NewRouter().StrictSlash(true)
	// myRouter.HandleFunc("/", hello).Methods(http.MethodPost)
	// log.Fatal(http.ListenAndServe(":10000", myRouter))
	var err error
	db, err = connectDB(connectionString)
	if err != nil {
		panic(err)
	}

	fmt.Println(autoMigrate())
	fmt.Println(createCompany("McDonalds"))
	fmt.Println(createSentiment("tweet1", 1, 0.1, "McDonalds"))
	fmt.Println(getSentiments("McDonalds", 10, 0))
}
