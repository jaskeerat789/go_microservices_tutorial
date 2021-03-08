package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

type Hello struct {
	l *log.Logger
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("hello world")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "Hello %s \n", d)
}
