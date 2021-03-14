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
	"log"
	"microservice_tutorial/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	l *log.Logger
}

type KeyProduct struct{}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// Responses:
// 	200: productResponse

// GetProducts returns the list of all products
func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProductData()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to read json", http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /products products registerProduct
// Adds a product to the Database
// Responses:
// 	201: noContent

// AddProduct add a Product to Products list
func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddToList(prod)

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
	err := data.UpdateProduct(id, prod)
	if err != nil {
		p.l.Printf("Error:%v", err)
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
	err := data.DeleteProduct(id)
	if err != nil {
		p.l.Printf("Error:%v", err)
		http.Error(rw, "Unable to delete Product", http.StatusInternalServerError)
		return
	}
}

// MiddlewareProductValidation decodes JSON
func (p Product) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Printf("Decodeing JSON error %v", err)
			http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Printf("Validation error %v", err)
			http.Error(rw, fmt.Sprintf("Wrong data passed: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
