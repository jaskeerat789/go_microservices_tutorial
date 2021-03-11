package handlers

import (
	"context"
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

func (p *Product) GetProductData(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProductData()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to read json", http.StatusInternalServerError)
		return
	}
}

func (p *Product) AddProductData(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddToList(prod)
	return
}

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
	return
}

func (p Product) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
