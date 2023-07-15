package models

import "time"

type Users struct {
	Email    string `json:"userEmail"`
	Password string `json:"password"`
}
type Store struct {
	Id       int
	Password []byte `json:"password"`
}
type Register struct {
	Name     string
	Type     string
	Email    string
	JoinedOn time.Time
}
type Item struct {
	Id      int       `json:"id"`
	TypeId  int       `json:"typeId" db:"type_id"`
	AddedBy int       `json:"AddedBy" db:"added_by"`
	Name    string    `json:"name" db:"item_name"`
	Price   int       `json:"price" db:"price"`
	AddedOn time.Time `db:"added_on"`
	Photos  []byte    `db:"photos"`
}
type ItemType struct {
	Type string `json:"type"`
}
type Pieces struct {
	Quantity int `json:"quantity"`
}
type CartItem struct {
	Id       int
	CartId   int    `json:"cartId"`
	ItemName string `json:"itemName"`
	ItemType string `json:"itemType"`
	ItemId   int    `json:"itemId"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Photos   []byte `db:"photos"`
}
type Checkout struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNo     int    `json:"phoneNo"`
	Address     string `json:"address"`
	ZipCode     int    `json:"zipCode"`
	City        string `json:"city"`
	Country     string `json:"country"`
	MoneyNumber int    `json:"moneyNumber"`
	MoneyPin    int    `json:"moneyPin"`
}
type Uploads struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Url  string `json:"url"`
}
type RegisterUser struct {
	Name     string `json:"userName"`
	Email    string `json:"userEmail"`
	Password string `json:"password"`
}