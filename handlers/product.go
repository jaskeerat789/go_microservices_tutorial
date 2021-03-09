package handlers

import (
	"log"
	"microservice_tutorial/data"
	"net/http"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProductData(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProductData(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProductData()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to read json", http.StatusInternalServerError)
	}
}
