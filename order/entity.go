package order

import (
	"time"

	"github.com/leekchan/accounting"
)

type Order struct {
	ID          	int		`gorm:"primaryKey;autoIncrement"`
	UserID     	 	int
	TotalPrice 		int
	Status      	string	`gorm:"size:50"`
	StatusPayment 	string	`gorm:"size:50"`
	Pengambilan		string	`gorm:"size:50"`
	CreatedAt   	time.Time	`gorm:"autoCreateTime"`
	Detail    		[]DetailOrder `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}


type DetailOrder struct {
	ID 			int		`gorm:"primaryKey;autoIncrement"`
	OrderID 	int		`gorm:"not null"`
	ProdukID	int
	ProdukName 	string	`gorm:"size:255"`
	Image 		string	`gorm:"size:255"`
	Price 		int
	Quantity 	int
	SubTotal 	int
	CreatedAt time.Time	`gorm:"autoCreateTime"`
}

func (c DetailOrder) PriceIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.Price)
}

func (c Order) TotPriceIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.TotalPrice)
}