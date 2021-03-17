package main

import (
	"context"
	"log"
	"microservice_tutorial/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jaskeerat789/gRPC-tutorial/protos/currency"
	"google.golang.org/grpc"
)

func main() {

	// initialize logger
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create new Currency Clients
	cc := currency.NewCurrencyClient(conn)

	// create handlers
	ph := handlers.NewProduct(l, cc)
	// create a new server mux and register the handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetProductById)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods("DELETE").Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	dh := middleware.Redoc(ops, nil)
	getRouter.Handle("/docs", dh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	goHandlers.CORS(goHandlers.AllowedHeaders([]string{"http://localhost:3000"}))

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
