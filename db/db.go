package db

import "golang.org/x/exp/slices"

var db []Product
var id_count int = 1

type Product struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Price   int    `json:"price"`
	ExpDate string `json:"exp_date"`
}

func Add(p Product) Product{
	p.Id = id_count
	id_count++
	db = append(db, p)
	return p
}
func Get() []Product {
	return db
}
func Putch(p Product, id int) Product{
	for i := range db {
		if id == db[i].Id {
			if p.Name!="" {
				db[i].Name=p.Name
			}
			if p.Desc!="" {
				db[i].Desc=p.Desc
			}
			if p.Price!=0 {
				db[i].Price=p.Price
			}
			if p.ExpDate!="" {
				db[i].ExpDate=p.ExpDate
			}
			return db[i]
		}
	}
	return Product{}                      
}

func Delete(id int) []Product {
	for i := range db {
		if id == db[i].Id {
			db=slices.Delete(db, i, i+1)
			return db
		}
	}
	return db
}