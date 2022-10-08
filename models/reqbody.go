package models


type Body struct {
	CustomerName string `json:"customerName"`
	Items []itemreq `json:"items"`
}

type itemreq struct {
	ItemID uint `json:"lineItemId"`
	ItemCode string `json:"itemCode"`
	Description string `json:"description"`
	Quantity int `json:"quantity"`
	OrderID uint
}