package types

import "time"

type Stock struct{
	Stkid        int      `gorm:"primary_key; unique_index; not null" json:"stkid"`
	Quant        int       `gorm:"not null" json:"quant"`
	Description  string    `gorm:"not null" json:"description"`
	Purchased_on time.Time `gorm:"not null" json:"purchased_on"`
	Initcost     float32   `gorm:"not null" json:"initcost"`
	U_id        int 		`gorm:"not null"  json:"uid"`
}
type Stocks []Stock

type PStock struct{
	Stk_id    int     `json:"stkid"`
	Stk_name  string  `json:"stk_name"`
	Cur_price float32 `json:"cur_price"`
	Pr_price  float32 `json:"pr_price"`
}
type Users struct{
	U_id   int    `gorm:"primary_key; unique_index; not null" json:"u_id"`
	U_name string `gorm:"primary_key; not null" json:"u_name"`
	U_pass string `gorm:"not null" json:"u_pass"`
	Amount float64 `gorm:"not null" json:"amount"`
	Email_id string `gorm:"not null" json:"email_id"`

}
type Market struct {
	Stkid int `json:"stkid"`
	Stk_name string `json:"stk_name"`
	Open_price float64 `json:"open_price"`
	Close_price float64 `json:"close_price"`
	High float64 `json:"high"`
	Low float64 `json:"low"`
	Bid string `json:"bid"`
	Ask string `json:"ask"`
	Volume int `json:"volume"`
}
type Markets []Market
type PStocks []PStock