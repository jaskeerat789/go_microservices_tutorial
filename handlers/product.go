// Package classification Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"microservice_tutorial/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jaskeerat789/gRPC-tutorial/protos/currency"
)

type Product struct {
	l          hclog.Logger
	productsDB *data.ProductsDB
}

type KeyProduct struct{}

func NewProduct(pdb *data.ProductsDB, l hclog.Logger) *Product {
	return &Product{l, pdb}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// Responses:
// 	200: productResponse

// GetProducts returns the list of all products
func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get all products")
	rw.Header().Add("Content-Type", "application/json")
	cur := r.URL.Query().Get("currency")

	lp, err := p.productsDB.GetProductData(cur)
	if err != nil {
		p.l.Error("Error:%v", err)
		http.Error(rw, "Unable to update Product", http.StatusInternalServerError)
		return
	}

	err = lp.ToJSON(rw)

	if err != nil {
		p.l.Error("Failed to serialize products:%v", err)
		http.Error(rw, "Failed to serialize products", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /products/{id} products listSingleProducts
// Returns a products
// Responses:
// 	200: singleProductResponse
func (p *Product) GetProductById(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get a product with id")
	rw.Header().Add("Content-Type", "application/json")
	cur := r.URL.Query().Get("currency")

	id, err := getProductId(r)
	if err != nil {
		p.l.Error("Id can not be retrieved: %v", err)
		http.Error(rw, "Id cannot be retrieved", http.StatusBadRequest)
		return
	}
	prod, err := p.productsDB.GetProductById(id, cur)
	if err != nil {
		p.l.Error("Product with id %v Does not exists: %v", id, err)
		http.Error(rw, fmt.Sprintf("Product with id as %v cannot be retrieved", id), http.StatusBadRequest)
		return
	}

	rr := &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_EUR),
		Destination: currency.Currencies(currency.Currencies_value["GBP"]),
	}
	resp, err := p.productsDB.Currency.GetRate(context.Background(), rr)

	if err != nil {
		p.l.Error("Failed to convert the currency", err)
		http.Error(rw, "Currency conversion", http.StatusInternalServerError)
		return
	}

	prod.Price = prod.Price * resp.Rate

	err = prod.ToJSON(rw)
	if err != nil {
		p.l.Error("Failed to serialize products:%v", err)
		http.Error(rw, "Unable to serialize the data", http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /products products registerProduct
// Adds a product to the Database
// Responses:
// 	201: noContent

// AddProduct add a Product to Products list
func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Add a product with id")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.productsDB.AddToList(prod)

}

// swagger:route PUT /products/{id} products updateProduct
// updates an existing product
// Responses:
// 	200: productResponse
//	500: InternalServerError

// GetProducts returns the list of all products
func (p *Product) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	id, _ := strconv.Atoi(vars["id"])
	err := p.productsDB.UpdateProduct(id, prod)
	if err != nil {
		p.l.Error("Error:%v", err)
		http.Error(rw, "Unable to update Product", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(201)

}

// swagger:route DELETE /products/{id} products removeProducts
// Removes an existing product from the DB
// Responses:
// 	201: noContent
//	500: InternalServerError

// DeleteProduct deletes a product from productList
func (p *Product) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := p.productsDB.DeleteProduct(id)
	if err != nil {
		p.l.Error("Error:%v", err)
		http.Error(rw, "Unable to delete Product", http.StatusInternalServerError)
		return
	}
}

func getProductId(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, error := strconv.Atoi(vars["id"])
	return id, error
}

// MiddlewareProductValidation decodes JSON
func (p Product) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Error("Decodeing JSON error %v", err)
			http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Error("Validation error %v", err)
			http.Error(rw, fmt.Sprintf("Wrong data passed: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
