package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Product struct {
	ID         int
	Name       string
	CategoryID int
	TypeID     int
	Model      string
	Price      float64
	Amount     int
	YearOfMade int
	Stores     []Store
}

type Category struct {
	ID   int
	Name string
}

type Type struct {
	ID   int
	Name string
}

type Store struct {
	ID       int
	Name     string
	Adresses []Adress
}

type Adress struct {
	ID       int
	District string
	Street   string
	StoreID  int
}

func main() {
	conString := `user=postgres password=1234 dbname=productdb sslmode=disable`
	db, err := sql.Open("postgres", conString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	products := []Product{
		{
			Name:       "iPhone 13",
			Model:      "13 PRO",
			Price:      1200,
			Amount:     10,
			CategoryID: 1,
			TypeID:     1,
			YearOfMade: 2010,
			Stores: []Store{
				{
					Name: "Malika bozor",
					Adresses: []Adress{
						{
							Street:   "Abu Saxiy",
							District: "Yunusobod",
						},
						{
							Street:   "Katta Qani",
							District: "Chilonzor",
						},
						{
							Street:   "Oybek street 10",
							District: "Uchtepa",
						},
					},
				},
				{
					Name: "Mediapark",
					Adresses: []Adress{
						{
							Street:   "Roxat Mediapark",
							District: "Yashnaobod",
						},
						{
							Street:   "Qoratosh mediapark",
							District: "Olmazor",
						},
						{
							Street:   "Chilonzor mediapark",
							District: "Chilonzor",
						},
					},
				},
			},
		},
		{
			Name:       "Stake",
			Model:      "Medium",
			Price:      100,
			Amount:     100,
			CategoryID: 2,
			TypeID:     2,
			YearOfMade: 2013,
			Stores: []Store{
				{
					Name: "Stake house",
					Adresses: []Adress{
						{
							District: "Chilonzor",
							Street:   "A.Qodiriy",
						},
						{
							District: "Mirzo Ulug`bek",
							Street:   "Qaysidir ko`cha",
						},
						{
							District: "Sergeli",
							Street:   "sergeli 10",
						},
						{
							District: "Farg`ona",
							Street:   "Jo`ydam",
						},
					},
				},
				{
					Name: "BurgerKing",
					Adresses: []Adress{
						{
							District: "Namangan",
							Street:   "To'da",
						},
						{
							District: "Pop",
							Street:   "Sang",
						},
						{
							District: "Uchkurgan",
							Street:   "Koson",
						},
					},
				},
			},
		},
		{
			Name:       "Chacman",
			Model:      "Man Y",
			Price:      125000,
			Amount:     35,
			CategoryID: 3,
			TypeID:     3,
			YearOfMade: 2008,
			Stores: []Store{
				{
					Name: "Chacman salon",
					Adresses: []Adress{
						{
							Street:   "Maxtumquli 12",
							District: "Chilonzor",
						},
					},
				},
			},
		},
		{
			Name:       "GUCCI jeans",
			Model:      "dior",
			Price:      355,
			Amount:     240,
			CategoryID: 4,
			TypeID:     4,
			YearOfMade: 2009,
			Stores: []Store{
				{
					Name: "Jeans shop",
					Adresses: []Adress{
						{
							Street:   "Uchtepa 2",
							District: "Uchtepa",
						},
						{
							Street:   "Qo'yliq 54",
							District: "Mirobod",
						},
					},
				},
			},
		},
	}
	for _, product := range products {
		var storeIDs []int
		for _, store := range product.Stores {
			var storeID int
			err := db.QueryRow(`insert into store (name) values ($1) returning id`, store.Name).Scan(&storeID)
			if err != nil {
				fmt.Println("error while inserting into store", err)
				return
			}
			storeIDs = append(storeIDs, storeID)
			for _, address := range store.Adresses {
				_, err := db.Exec(`insert into adresses (street, district, store_id) values ($1, $2, $3)`, address.Street, address.District, storeID)
				if err != nil {
					fmt.Println("error while inserting into addresses", err)
					return
				}
			}
		}

		err = db.QueryRow(`insert into products (name, model, price, amount, type_id, category_id) values ($1, $2, $3, $4, $5, $6) returning id`, product.Name, product.Model, product.Price, product.Amount, product.TypeID, product.CategoryID).Scan(&product.ID)
		if err != nil {
			fmt.Println("error while inserting into products", err)
			return
		}
		for _, strID := range storeIDs {
			_, err = db.Exec(`insert into stores_products (store_id, product_id) values ($1, $2)`, strID, product.ID)
			if err != nil {
				fmt.Println("error while inserting into store products", err)
				return
			}
		}
	}
}
