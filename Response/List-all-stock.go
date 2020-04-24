package responses

import (
	"GO-INVEST/types"
	"net/http"
)

type StocksResponse struct {
	*types.Stocks
}

func ListStocks(expenses *types.Stocks) *StocksResponse{
	return &StocksResponse{ expenses}

}

func (e *StocksResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return nil
}

