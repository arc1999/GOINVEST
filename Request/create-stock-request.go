package requests

import (
	"GO-INVEST/types"
	"errors"
	"net/http"
)

type BuyStockRequest struct {
	*types.Stock
}

func (c *BuyStockRequest) Bind(r *http.Request) error {
	if c.Description == "" {
		return errors.New("description is either empty or invalid")
	}

	if c.Initcost==0 {
		return errors.New("amount is either empty or invalid")
	}

	if c.Quant == 0 {
		return errors.New("Quant is either zero or invalid")
	}

	return nil
}