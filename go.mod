module microservice_tutorial

go 1.16

require (
	github.com/go-openapi/runtime v0.19.26
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v0.15.0
	github.com/jaskeerat789/gRPC-tutorial v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.36.0
)

replace github.com/jaskeerat789/gRPC-tutorial => ../gRPC
