module microservice_tutorial

go 1.16

require (
	github.com/go-openapi/runtime v0.19.26
	github.com/go-playground/validator/v10 v10.4.1
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/jaskeerat789/gRPC-tutorial v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	google.golang.org/grpc v1.36.0
)

replace github.com/jaskeerat789/gRPC-tutorial => ../gRPC
