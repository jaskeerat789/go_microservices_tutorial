package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/jaskeerat789/gRPC-tutorial/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"sku"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
}

type Products []*Product

type ProductsDB struct {
	Currency currency.CurrencyClient
	Log      hclog.Logger
	rates    map[string]float64
	client   currency.Currency_SubscribeRatesClient
}

func NewProductDB(c currency.CurrencyClient, l hclog.Logger) *ProductsDB {
	pd := &ProductsDB{c, l, make(map[string]float64), nil}
	go pd.handleUpdates()
	return pd
}

func (pdb *ProductsDB) GetProductData(Currency string) (Products, error) {
	if Currency == "" {
		return productList, nil
	}

	rate, err := pdb.getRate(Currency)
	if err != nil {
		pdb.Log.Error("Error: %v", err)
		return nil, err
	}

	pr := Products{}
	for _, product := range productList {
		np := *product
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}
	return pr, nil
}

func (pdb *ProductsDB) GetProductById(id int, Currency string) (*Product, error) {
	pos, err := getPos(id)
	if err != nil {
		return &Product{}, fmt.Errorf("Product with id as %v not found", id)
	}

	if Currency == "" {
		return productList[pos], nil
	}

	rate, err := pdb.getRate(Currency)
	if err != nil {
		pdb.Log.Error("Error: %v", err)
		return nil, err
	}

	pr := *productList[pos]
	pr.Price = pr.Price * rate
	return &pr, nil

}

func (pdb *ProductsDB) AddToList(p *Product) {
	p.ID = getId()
	productList = append(productList, p)
}

func (pdb *ProductsDB) UpdateProduct(id int, p *Product) error {
	pos, err := getPos(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p
	return nil
}

func (pdb *ProductsDB) DeleteProduct(id int) error {
	pos, err := getPos(id)
	if err != nil {
		return ErrorProductNotFound
	}
	productList[pos] = productList[len(productList)-1]
	productList = productList[:len(productList)-1]
	return nil
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

var ErrorProductNotFound = fmt.Errorf("Product not found")

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

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {

	if r, ok := p.rates[destination]; ok {
		return r, nil
	}

	rr := &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_EUR),
		Destination: currency.Currencies(currency.Currencies_value[destination]),
	}
	resp, err := p.Currency.GetRate(context.Background(), rr)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*currency.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("unable to get rate from currency server, destination and base currency can not be same, base: %s, destination: %s", md.Base.String(), md.Destination.String())

			}
			return -1, fmt.Errorf("unable to get rate from currency server, base: %s, destination: %s", md.Base.String(), md.Destination.String())
		}

	}

	p.client.Send(rr)
	p.rates[destination] = resp.Rate

	return resp.Rate, err
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.Currency.SubscribeRates(context.Background())
	if err != nil {
		p.Log.Error("Unable to subscribe for rates", "error", err)
		return
	}

	p.client = sub

	for {
		rr, err := sub.Recv()
		p.Log.Info("Received updated rate from server", "dest", rr.GetDestination().String())
		if err != nil {
			p.Log.Error("Error receiving message", "error", err)
			return
		}
		p.rates[rr.Destination.String()] = rr.Rate
	}

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
