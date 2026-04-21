package db

import (
	"errors"
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/exp/slices"
)

var db []Product
var idCount int = 1
var ErrNotFound = errors.New("product not found")
type Product struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Price   int    `json:"price"`
	ExpDate string `json:"exp_date"`
}

func Add(p Product) Product{
	p.Id = idCount
	idCount++
	db = append(db, p)
	log.Infof("Product created: id=%d", p.Id)
	return p
}
func Get() []Product {
	log.Debugf("Returning all products")
	return db
}
func Patch(p Product, id int) (Product,error){
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
			log.Infof("Product with id: %d updated", id)
			return db[i], nil
		}
	}
	log.Warnf("DB: product not found id: %d",id)
	return Product{}, ErrNotFound                 
}

func Delete(id int) error {
	for i := range db {
		if id == db[i].Id {
			db=slices.Delete(db, i, i+1)
			log.Infof("Product with id: %d deleted", id)
			return nil
		}
	}
	return ErrNotFound 
}