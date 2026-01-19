package handler

import (
	"fmt"
	"net/http"
	"toko/order"
	"toko/produk"
	"toko/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type adminHandler struct {
	userService user.Service
	produkService produk.Service
	orderService order.Service
}

func NewAdminHandler(userService user.Service, produkService produk.Service,orderService order.Service) *adminHandler{
	return &adminHandler{userService, produkService, orderService}
}

func(h *adminHandler) IndexAdmin(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	admin := session.Get("userRole")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Akses Ditolak",
		})
		return
	}

	user, err :=h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "index 1",
		})
		return
	}
	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	users, err := h.userService.GetAllUser()
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html",gin.H{
			"code":    500,
			"message":"Gagal memuat users",
			"title":   "index 2",
		})
		return
	}
	
	c.HTML(http.StatusOK,"indexAdmin.html", gin.H{
		"user":user,
		"avatar": avatar,
		"isLoggedIn": true,
		"users":users,
	})
}

func(h *adminHandler) ShouProduk(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	admin := session.Get("userRole")
	successMsg := session.Get("success")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Akses Ditolak",
		})
		return
	}

	user, err :=h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "index 1",
		})
		return
	}
	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	category := c.Query("category")
	name := c.Query("name")

	var produks []produk.Produk
	if name != "" {
		produks, err = h.produkService.GetByName(name)
	} else if category != "" {
		produks, err = h.produkService.GetByCategory(category)
	} else {
		produks, err = h.produkService.GetAllProduk()
	}


	
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html",gin.H{
			"code":    500,
			"message":"Gagal memuat produk",
			"title":   "index 2",
		})
		return
	}

	if successMsg != nil {
		session.Delete("success")
		session.Save()
	}

	c.HTML(http.StatusOK,"produkAdmin.html", gin.H{
		"user":user,
		"avatar": avatar,
		"isLoggedIn": true,
		"produks": produks,
		"keyword": category,
		"namekey": name,
		"success":successMsg,
	})
}

func(h *adminHandler) AddGetProduk(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	admin := session.Get("userRole")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Akses Ditolak",
		})
		return
	}

	c.HTML(http.StatusOK,"addProduk.html", nil)
}
func(h *adminHandler) AddPostProduk(c *gin.Context){
	var input produk.InputAddProduk

	session := sessions.Default(c)
	userID:= session.Get("userID")
	admin := session.Get("userRole")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Password Akses Ditolak",
		})
		return
	}

	err := c.ShouldBind(&input)
	if err !=nil {
		c.HTML(http.StatusBadRequest, "addProduk.html", gin.H{
			"error":      "Semua kolom harus diisi.",
			"title":      "add produk",
			"isLoggedIn": true,
		})
		return
	}

	file , err := c.FormFile("image")
	if err != nil {
        c.HTML(http.StatusBadRequest, "addProduk.html", gin.H{
            "error":      "Gambar harus dipilih.",
            "title":      "add produk",
            "isLoggedIn": true,
        })
        return
    }

	lastID, err:= h.produkService.GetlastProdukID()
	if err != nil {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"error":      "400",
			"title":      "tambah produk gagal",
			"isLoggedIn": true,
		})
		return
	}

	idPlus := lastID + 1
	
	imagePath :=fmt.Sprintf("images/produk/%d-%s",idPlus,file.Filename)
	err = c.SaveUploadedFile(file,imagePath)
	 if err != nil {
        c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
            "error":      "Gagal menyimpan gambar.",
            "title":      "add produk",
            "isLoggedIn": true,
        })
        return
    }

	_, err = h.produkService.AddProduk(input, imagePath)
	if err !=nil {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"error":      "400",
			"title":      "tambah produk gagal",
			"isLoggedIn": true,
		})
		return
	}
	session.Set("success", "Produk berhasil di tambahkan")
	session.Save()

	c.Redirect(http.StatusFound, "/produkAdmin")
}

func(h *adminHandler) UpdateGetProduk(c *gin.Context){
	session := sessions.Default(c)
	userID:= session.Get("userID")
	admin := session.Get("userRole")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Password Akses Ditolak",
		})
		return
	}

	var uri produk.GetProdukUri
	err := c.ShouldBindUri(&uri) 
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID barang tidak valid",
			"title":   "input param id produk error",
		})
		return
	}
	getproduk , err := h.produkService.GetProdukById(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat produk",
			"title":   "get produk error",
		})
		return
	}

	c.HTML(http.StatusOK,"updateProduk.html", gin.H{
		"produk":getproduk,
		"isLoggedIn": true,
	})

}

func(h *adminHandler) UpdatePostProduk(c *gin.Context){
	session := sessions.Default(c)
	userID:= session.Get("userID")
	admin := session.Get("userRole")
	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Password Akses Ditolak",
		})
		return
	}

	var uri produk.GetProdukUri
	err := c.ShouldBindUri(&uri) 
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID barang tidak valid",
			"title":   "input param id produk error",
		})
		return
	}

	var input produk.InputUpdateProduk
	err = c.ShouldBind(&input)
	if err !=nil {
		// Jika error-nya bukan karena field kosong (EOF), tampilkan error
		c.HTML(http.StatusBadRequest, "updateProduk.html", gin.H{
			"error":      "Terjadi kesalahan saat memproses data.",
			"title":      "Update Produk",
			"isLoggedIn": true,
		})
		return
	}

	file , err := c.FormFile("image")
	if err == nil {
       idPlus := uri
	
		imagePath :=fmt.Sprintf("images/produk/%d-%s", idPlus ,file.Filename)
		err = c.SaveUploadedFile(file,imagePath)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
				"error":      "Gagal menyimpan gambar.",
				"title":      "add produk",
				"isLoggedIn": true,
			})
			return
		}
    }
	
	imagePath := ""
	
	_, err = h.produkService.UpdateProduk(uri,input,imagePath)
	if err !=nil {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"error":      "400",
			"title":      "Update produk gagal",
			"isLoggedIn": true,
		})
		return
	}
	session.Set("success", "Produk berhasil di update")
	session.Save()

	c.Redirect(http.StatusFound, "/produkAdmin")

}

func(h *adminHandler) DeletePostProduk(c *gin.Context){
	session := sessions.Default(c)
	userID:= session.Get("userID")
	admin := session.Get("userRole")
	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Password Akses Ditolak",
		})
		return
	}

	var uri produk.GetProdukUri
	err := c.ShouldBindUri(&uri) 
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID barang tidak valid",
			"title":   "input param id produk error",
		})
		return
	}

	err = h.produkService.DeleteProduk(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "gagal menghapus barang",
			"title":   "input param id produk error",
		})
		return
	}

	session.Set("success", "Produk berhasil dihapus")
	session.Save()

	c.Redirect(http.StatusFound, "/produkAdmin")


}

func(h *adminHandler) GetLihatTransaksi(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	admin := session.Get("userRole")
	successMsg := session.Get("success")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Akses Ditolak",
		})
		return
	}

	user, err :=h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "index 1",
		})
		return
	}
	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	var order []order.Order
	status := c.Query("status")
	if status != "" {
		order, err = h.orderService.GetByStatus(status)
	}else{
		order, err = h.orderService.GetAllOrder()
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat transaksi",
			"title":   "transaction",
		})
		return
	}

	
	if successMsg != nil {
		session.Delete("success")
		session.Save()
	}

	c.HTML(http.StatusOK,"lihatTransaksi.html",gin.H{
		"user":user,
		"avatar": avatar,
		"order":order,
		"statuskey":status,
		"isLoggedIn": true,
		"success":successMsg,
	})

}

func(h *adminHandler) PostUpdateStatus(c *gin.Context){
	session := sessions.Default(c)
	
	var uri order.UriInputOrderID
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID order tidak valid",
			"title":   "input param id order error",
		})
		return
	}

	var input order.UpdateStatus
	err = c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "input order error ",
			"title":   "update status",
		})
		return
	}

	err = h.orderService.UpdateStatusOrder(uri,input.Status)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "update status gagal",
			"title":   "update  order status gagal",
		})
		return
	}

	session.Set("success", "Produk berhasil di update")
	session.Save()

	c.Redirect(http.StatusFound,"/produkAdmin/Transaksi")
}

func(h *adminHandler) PostUpdateStatusPay(c *gin.Context){
	session := sessions.Default(c)
	var uri order.UriInputOrderID
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID order tidak valid",
			"title":   "input param id order error",
		})
		return
	}

	var input order.UpdateStatusPay
	err = c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "input order error ",
			"title":   "update status",
		})
		return
	}

	err = h.orderService.UpdateStatusPayOrder(uri,input.Statuspayment)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "update status payment gagal",
			"title":   "update order status payment gagal",
		})
		return
	}

	session.Set("success", "Produk berhasil di update")
	session.Save()

	c.Redirect(http.StatusFound,"/produkAdmin/Transaksi")
}

func (h *adminHandler) GetDetailTransaction(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	admin := session.Get("userRole")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Akses Ditolak",
		})
		return
	}

	user, err :=h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "index 1",
		})
		return
	}
	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	var uri order.UriInputOrderID
	err = c.ShouldBindUri(&uri)
	if err !=nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID order tidak valid",
			"title":   "input param id order error",
		})
		return
	}

	getDetail , err := h.orderService.GetDetailOrder(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat detail order",
			"title":   "get produk error",
		})
		return
	}

	order , err := h.orderService.GetOrderById(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat detail order",
			"title":   "get produk error",
		})
		return
	}

	getPelanggan, err := h.userService.GetUserByID(order.UserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat",
			"title":   "get produk error",
		})
		return
	}

	

	c.HTML(http.StatusOK,"adminDetailTransaksi.html", gin.H{
		"user":user,
		"nama":getPelanggan,
		"avatar": avatar,
		"tgl" : order.CreatedAt.Format("2006-01-02"),
		"isLoggedIn": true,
		"Detail":getDetail,
		"total" : order.TotPriceIDR(),
	})
}

func (h *adminHandler) GetDetailUser(c *gin.Context){
	session := sessions.Default(c)
	userID:= session.Get("userID")
	admin := session.Get("userRole")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Password Akses Ditolak",
		})
		return
	}

	ser, err :=h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "index 1",
		})
		return
	}
	avatar := ser.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	var uri user.UserUri
	err = c.ShouldBindUri(&uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID user tidak valid",
			"title":   "input param id produk error",
		})
		return
	}

	getuser, err := h.userService.GetDetailUser(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user detail",
			"title":   "get user Detail error",
		})
		return
	}
	avat := getuser.Avatar
	if  avat ==""{
		avat = "images/avatar/avatar.jpg"
	}
	

	c.HTML(http.StatusOK,"userDetail.html", gin.H{
		"users":getuser,
		"avat" : avat,
		"isLoggedIn": true,
		"user":ser,
		"avatar": avatar,
	})
}

func ( h *adminHandler) DetailProduk(c *gin.Context){
	session := sessions.Default(c)
	userID:= session.Get("userID")
	admin := session.Get("userRole")

	if admin == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusBadRequest, "error2.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Password Akses Ditolak",
		})
		return
	}

	user, err :=h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "index 1",
		})
		return
	}
	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	var uri produk.GetProdukUri
	err = c.ShouldBindUri(&uri) 
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error2.html", gin.H{
			"code":    500,
			"message": "ID barang tidak valid",
			"title":   "input param id produk error",
		})
		return
	}
	getproduk , err := h.produkService.GetProdukById(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error2.html", gin.H{
			"code":    500,
			"message":"Gagal memuat produk",
			"title":   "get produk error",
		})
		return
	}

	c.HTML(http.StatusOK,"produkDetail2.html", gin.H{
		"produk":getproduk,
		"user":user,
		"avatar": avatar,
		"isLoggedIn": true,
	})

}