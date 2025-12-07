package order

type CartItem struct {
	ProductID   int
	ProductName string
	ProdukImage string
	Quantity    int
	Price       int
	SubTotal    int
}

type Peng struct {
	Pengambilan string `form:"pengambilan" binding:"required"`
}

type UriInputOrderID struct {
	ID int `uri:"id" binding:"required"`
}

type UpdateStatus struct {
	Status string `form:"status"`
}

type UpdateStatusPay struct {
	Statuspayment string `form:"status_payment"`
}