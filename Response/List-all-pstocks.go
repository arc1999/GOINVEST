package responses

import (
	"GO-INVEST/types"
	"net/http"
)

type PStocksResponse struct {
	*types.PStocks
}

func ListPStocks(expenses *types.PStocks) *PStocksResponse{
	return &PStocksResponse{ expenses}

}

func (e *PStocksResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return nil
}