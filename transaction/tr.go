package main

import (
	"database/sql"
	"fmt"
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

type ProductResponse struct {
	ID       int
	Name     string
	Category Category
	Type     Type
	Model    string
	Price    float64
	Amount   int
	Stores   []Store
	
}

func main() {
	conString := `user=postgres password=1234 dbname=productdb sslmode=disable`
	db, err := sql.Open("postgres", conString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	products := []ProductResponse{}

	rows, err := db.Query(`select
		id,
		name,
		type_id,
		category_id,
		price,
		amount,
		model from products`)
	if err != nil {
		fmt.Println("error while getting all products", err)
		return
	}

	for rows.Next() {
		product := ProductResponse{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Type.ID,
			&product.Category.ID,
			&product.Price,
			&product.Amount,
			&product.Model,
		)
		if err != nil {
			fmt.Println("error while scanning the product", err)
			return
		}

		err = db.QueryRow(`select name from categories where id = $1`, 
		product.Category.ID).Scan(&product.Category.Name)
		if err != nil {
			fmt.Println("error while getting product category", err)
			return
		}
		err = db.QueryRow(`select name from type where id = $1`,
		product.Type.ID).Scan(&product.Type.Name)	
		if err != nil {
			fmt.Println("error while getting product type", err)
			return
		}
		storeRows, err := db.Query(`select id, name from store s join stores_products sp on s.id=sp.store_id where sp.product_id = $1`, product.ID)
		if err != nil {
			fmt.Println("error while getting product stores", err)
			return
		}

		for storeRows.Next() {
			store := Store{}
			err := storeRows.Scan(
				&store.ID,
				&store.Name,
			)
			if err != nil {
				fmt.Println("error while getting product store info", err)
				return
			}

			addressRows, err := db.Query(`select id, district, street from adresses where store_id = $1`, store.ID)
			if err != nil {
				fmt.Println("error while getting store addresses", err)
				return
			}

			for addressRows.Next() {
				address := Adress{}
				err := addressRows.Scan(
					&address.ID,
					&address.District,
					&address.Street,
				)
				if err != nil {
					fmt.Println("error while getting store address info", err)
					return
				}

				store.Adresses = append(store.Adresses, address)
			}

			product.Stores = append(product.Stores, store)
		}

		products = append(products, product)
	}
	fmt.Println(products)
}
