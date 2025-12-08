package produk

import (
	"time"

	"github.com/leekchan/accounting"
)

type Produk struct {
	ID          int		`gorm:"primaryKey;autoIncrement;unsigned"`
	Name        string	`gorm:"size:150"`
	ImageName 	string	`gorm:"size:255"`
	Description string	`gorm:"type:text"`
	Category 	string	`gorm:"size:100"`
	Stock       int
	Price       int
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (c Produk) PriceIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.Price)
}