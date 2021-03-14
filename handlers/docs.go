package handlers

import "microservice_tutorial/data"

// A List of Products returns in the response
// swagger:response productResponse
type productResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters removeProducts updateProduct
type productIDParameterWrapper struct {
	// The id of the product to be removed from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

// No content is returned
// swagger:response noContent
type productNoContent struct {
}

// Internal Server Error response
// swagger:response  InternalServerError
type productInternalServerError struct {
	// Error response
	// in: body
	Body []errorResponse
}

//swagger:model
type errorResponse struct {
	// Error code
	// in: body
	Err int `json:"err"`
	// Error message
	// in: body
	Message string `json:"message"`
}
