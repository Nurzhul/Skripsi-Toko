package produk

type InputAddProduk struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
	Category    string `form:"category" binding:"required"`
	Stock       int    `form:"stock" binding:"required"`
	Price       int    `form:"price" binding:"required"`
}

type GetProdukUri struct {
	ID int `uri:"id" binding:"required"`
}

type InputUpdateProduk struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Category    string `form:"category"`
	Stock       int    `form:"stock"`
	Price       int    `form:"price"`
}