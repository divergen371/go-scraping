package main

import (
	"fmt"
)

func main() {
	baseURL := "http://localhost:5001/"
	resp, err := fetch(baseURL)

	if err != nil {
		panic(err)
	}

	items, err := parseList(resp)
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		fmt.Println(item)
	}

	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	err = migrateDB(db)
	if err != nil {
		panic(err)
	}

	err = createLatestItems(items, db)
	if err != nil {
		panic(err)
	}

	err = updateItemMaster(db)
	if err != nil {
		panic(err)
	}
}
