package responses

import (
	"GO-INVEST/types"
	"net/http"
)

type ListStockResponse struct {
	*types.Stock
}


func (ListStockResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListStock(exp *types.Stock) *ListStockResponse {
	resp := &ListStockResponse{Stock: exp}
	return resp
}