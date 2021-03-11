package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

type Products []*Product

func GetProductData() Products {
	return productList
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func AddToList(p *Product) {
	p.ID = getId()
	productList = append(productList, p)
}

var ErrorProductNotFound = fmt.Errorf("Product not found")

func UpdateProduct(id int, p *Product) error {
	pos, err := getPos(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p
	return nil
}

func getPos(id int) (int, error) {
	for i, prod := range productList {
		if prod.ID == id {
			return i, nil
		}
	}
	return -1, ErrorProductNotFound
}

func getId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedAt:   time.Now().UTC().String(),
		DeletedAt:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk ",
		Price:       1.99,
		SKU:         "fgh123",
		CreatedAt:   time.Now().UTC().String(),
		DeletedAt:   time.Now().UTC().String(),
	},
}
