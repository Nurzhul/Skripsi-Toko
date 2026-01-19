package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"toko/helper"
	"toko/order"
	"toko/produk"
	"toko/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	userService user.Service
	produkService produk.Service
	orderService order.Service
}

func NewUserHandler(userService user.Service,  produkService produk.Service, orderService order.Service) *userHandler{
	return &userHandler{userService, produkService,orderService}
}


func (h *userHandler) Index(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("userID")

    var (
        user       user.User
        avatar     string
        isLoggedIn bool
    )

    // Jika user belum login → tetap lanjut, tapi gunakan data default
    if userID != nil {
        // Cek apakah tipe data sesuai
        id, ok := userID.(int)
        if ok {
            // Ambil data user
            usr, err := h.userService.GetUserByID(id)
            if err == nil {
                user = usr
                isLoggedIn = true
            }
        }
    }

    // Set avatar default
    if user.Avatar != "" {
        avatar = user.Avatar
    } else {
        avatar = "/images/avatar/avatar.jpg"
    }

    c.HTML(http.StatusOK, "index.html", gin.H{
        "user":       user,
        "avatar":     avatar,
        "isLoggedIn": isLoggedIn,
    })
}

func (h *userHandler) New(c *gin.Context){
	c.HTML(http.StatusOK,"register.html",nil)
}

func (h *userHandler) Regis(c *gin.Context){
	var input user.InputRegister

	err := c.ShouldBind(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"errors": ve,
			})
		}else {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": "Invalid input",
			})
		}
		return
	}
	// log.Println("Sebelum RegisterUser dipanggil")
	_, err = h.userService.RegisterUser(input)
	if err != nil {
		c.HTML(http.StatusBadRequest,"register.html",gin.H{
			"emailError": err.Error(), // kirim pesan error ke halaman
		})
		return
	}

	// log.Println("Setelah RegisterUser dipanggil")
	c.Redirect(http.StatusFound, "/login?register_success=true")
}

func (h *userHandler) Newpassword(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}

	c.HTML(http.StatusOK,"update_pas.html", nil)
}

func (h *userHandler) UpdatePass(c *gin.Context){
	session := sessions.Default(c)
	userID:= session.Get("userID")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Password Akses Ditolak",
		})
		return
	}

	var input user.ChangePassword

	
	err := c.ShouldBind(&input)
	if err !=nil {
		c.HTML(http.StatusBadRequest, "update_pas.html", gin.H{
			"error":      "panjang password minimal 8 karakter.",
			"title":      "Ganti Password",
			"isLoggedIn": true,
		})
		return
	}

	if input.NewPassword != input.CpmfirmPassword{
		c.HTML(http.StatusBadRequest, "update_pas.html", gin.H{
			"error":      "Konfirmasi password tidak sama.",
			"title":      "Ganti Password",
			"isLoggedIn": true,
		})
		return
	}

	

	_, err = h.userService.Updatepassword(id,input)
	if err != nil {
		c.HTML(http.StatusBadRequest, "update_pas.html", gin.H{
			"error":      err.Error(),
			"title":      "Ganti Password",
			"isLoggedIn": true,
		})
		return
	}
	session.Set("success", "Password berhasil diubah.")
	session.Save()

	c.Redirect(http.StatusFound, "/profile")
}

func (h *userHandler) Profile(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	successMsg := session.Get("success")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Profile 1",
		})
		return
	}

	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	if successMsg != nil {
		session.Delete("success")
		session.Save()
	}


	c.HTML(http.StatusOK,"profile.html", gin.H{
		"user":user,
		"success":successMsg,
		"avatar": avatar,
		"isLoggedIn": true,
	})
}

func (h *userHandler) Addres(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}
	c.HTML(http.StatusOK,"addres.html", nil)
}

func (h *userHandler) NewAddres(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Addres Akses Ditolak",
		})
		return
	}

	var input user.InputAddres
	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusBadRequest, "addres.html", gin.H{
			"error":      "Semua kolom harus diisi",
			"title":      "Lengkapi Form",
			"isLoggedIn": true,
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err !=nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Addres 1",
		})
		return		
	}

	if user.Addres.ID == 0 {
		_, err = h.userService.SaveAddres(id, input)
		if err != nil{
			c.HTML(http.StatusInternalServerError,"error.html", gin.H{
				"code":    500,
				"message":"Server mengalami kesalahan",
				"title":   "Addres 2",
			})
			return		
		}
		session.Set("success", "Alamat berhasil ditambahkan.")
		session.Save()
	}else{
		_, err = h.userService.UpdateAddres(id, input)
		if err != nil {
			c.HTML(http.StatusInternalServerError,"error.html", gin.H{
				"code":    500,
				"message":"Server mengalami kesalahan",
				"title":   "Addres 3",
			})
		}
		session.Set("success", "Alamat berhasil diubah.")
		session.Save()
	}

	c.Redirect(http.StatusFound, "/profile")

}
func (h *userHandler) SaveGetAvatar(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Profile 1",
		})
		return
	}

	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	c.HTML(http.StatusOK,"PhotoProfile.html",gin.H{
		"user":user,
		"avatar": avatar,
		"isLoggedIn": true,
	})
}
func (h *userHandler) SavePostAvatar(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"terjadi kesalahan ketika memuat user",
			"title":   "Get id",
		})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusBadRequest, "PhotoProfile.html", gin.H{
			"code":      400,
			"title":     "Gambar harus dipilih",
			"message":   "Tidak ada foto yang dikirim",
			"isLoggedIn": true,
		})
		return
	}

	imagePath :=fmt.Sprintf("images/avatar/%d-%s",id,file.Filename)
	err = c.SaveUploadedFile(file,imagePath)
	 if err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "code":      500,
            "title":      "upload avatar gagal",
			"message": "terjadi kesalahan ketika mengupload avatar 2",
            "isLoggedIn": true,
        })
        return
    }

	
	_, err = h.userService.SaveAvatar(user.ID, imagePath)
	if err !=nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":      500,
			"title":      "save avatar gagal",
			"message": "terjadi kesalahan ketika menyimpan avatar",
			"isLoggedIn": true,
		})
		return
	}

	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}
	session.Set("success", "Avatar berhasil di ubah")
	session.Save()

	c.Redirect(http.StatusFound,"/profile")

}

func (h *userHandler) Produk (c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")

	var (
		user      user.User // misalnya struct user kamu
		avatar    string
		isLoggedIn bool
	)

	
	if userID != nil {
		id, ok := userID.(int)
		if !ok {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"code":    500,
				"message": "ID pengguna tidak valid",
				"title":   "Akses Ditolak",
			})
			return
		}

		var err error
		user, err =h.userService.GetUserByID(id)
		if err != nil {
			c.HTML(http.StatusInternalServerError,"error.html", gin.H{
				"code":    500,
				"message":"Gagal memuat user",
				"title":   "index 1",
			})
			return
		}
		isLoggedIn = true
		avatar = user.Avatar
		if  avatar ==""{
			avatar = "/images/avatar/avatar.jpg"
		}
	}
	

	category := c.Query("category")
	name := c.Query("name")

	var err error
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
	
	c.HTML(http.StatusOK,"produk.html", gin.H{
		"user":user,
		"avatar": avatar,
		"isLoggedIn": isLoggedIn,
		"produks": produks,
		"keyword": category,
		"namekey": name,
	})
}

func (h *userHandler) DetailProduk (c *gin.Context){
	Session := sessions.Default(c)
	userID := Session.Get("userID")
	
	var (
		user      user.User // misalnya struct user kamu
		avatar    string
		isLoggedIn bool
		err        error
	)

	
	if userID != nil {
		id, ok := userID.(int)
		if !ok {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"code":    500,
				"message": "ID pengguna tidak valid",
				"title":   "Akses Ditolak",
			})
			return
		}

		
		user, err =h.userService.GetUserByID(id)
		if err != nil {
			c.HTML(http.StatusInternalServerError,"error.html", gin.H{
				"code":    500,
				"message":"Gagal memuat user",
				"title":   "index 1",
			})
			return
		}
		isLoggedIn = true
		avatar = user.Avatar
		if  avatar ==""{
			avatar = "images/avatar/avatar.jpg"
		}
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

	c.HTML(http.StatusOK,"produkDetail.html", gin.H{
		"user":user,
		"avatar": avatar,
		"isLoggedIn": isLoggedIn,
		"produk":getproduk,
		"ID":getproduk.ID,
	})

}

func (h *userHandler) AddToCart (c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Addres Akses Ditolak",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err !=nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Addres 1",
		})
		return		
	}

	if user.Addres.ID == 0 {
		session.Set("success", "isikan alamat dahulu")
		session.Save()
		c.Redirect(http.StatusFound, "/profile")
	}

	var input struct {
        ProductID int `json:"product_id" binding:"required"`
        Quantity  int `json:"quantity" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("inputJson")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
        return
    }

    product, err := h.produkService.GetProdukById(produk.GetProdukUri{
		ID: input.ProductID,
	})
	
    if err != nil {
		log.Println("produk id service")
        c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
        return
    }

   
     var items []order.CartItem

    cartData := session.Get("cart")
    if cartData != nil {
        // decode JSON ke slice struct
        json.Unmarshal([]byte(cartData.(string)), &items)
    }

    found := false
    for i, item := range items {
        if item.ProductID == input.ProductID {
            items[i].Quantity += input.Quantity
            found = true
            break
        }
    }

    if !found {
        items = append(items, order.CartItem{
            ProductID:   product.ID,
            ProductName: product.Name,
			ProdukImage: product.ImageName,
            Quantity:    input.Quantity,
            Price:       product.Price,
			SubTotal:input.Quantity*product.Price,
        })
    }

    cartJSON, _ := json.Marshal(items)
    session.Set("cart", string(cartJSON))
    session.Save()

	// Hitung total item di keranjang
    totalCount := 0
    for _, item := range items {
        totalCount += item.Quantity
    }


    c.JSON(http.StatusOK, gin.H{
		"message": "Produk ditambahkan ke keranjang",
		"cart_count": totalCount,
	})
}

func (h *userHandler) GetCartCount(c *gin.Context) {
    session := sessions.Default(c)
    cartData := session.Get("cart")

    count := 0

    if cartData != nil {
        var items []order.CartItem

        // decode JSON
        if err := json.Unmarshal([]byte(cartData.(string)), &items); err == nil {
            for _, item := range items {
                count += item.Quantity
            }
        }
    }

    c.JSON(http.StatusOK, gin.H{"cart_count": count})
}

func (h *userHandler) ViewCart(c *gin.Context) {
    session := sessions.Default(c)
    cartData := session.Get("cart")
    userID := session.Get("userID")

    if userID == nil {
        c.Redirect(http.StatusFound, "/login")
        return
    }

    id, ok := userID.(int)
    if !ok {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "code":    500,
            "message": "ID pengguna tidak valid",
            "title":   "Addres Akses Ditolak",
        })
        return
    }

    user, err := h.userService.GetUserByID(id)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "code":    500,
            "message": "Gagal memuat user",
            "title":   "Profile 1",
        })
        return
    }
	
	if user.Addres.ID == 0 {
		session.Set("success", "isikan alamat dahulu")
		session.Save()
		c.Redirect(http.StatusFound, "/profile")
	}

    avatar := user.Avatar
    if avatar == "" {
        avatar = "/images/avatar/avatar.jpg"
    }

    // Decode cart
    var items []order.CartItem
    total := 0

    if cartData != nil {
        // Decode JSON
        if err := json.Unmarshal([]byte(cartData.(string)), &items); err != nil {
            log.Println("ERROR decode cart:", err)
        } else {
            for _, item := range items {
                total += item.Quantity * item.Price
            }
        }
    }


    c.HTML(http.StatusOK, "keranjang.html", gin.H{
        "CartItems":  items,
        "Total":      total,
        "user":       user,
        "avatar":     avatar,
        "isLoggedIn": true,
    })
}

func (h *userHandler) Checkout(c *gin.Context) {
    session := sessions.Default(c)
    cartData := session.Get("cart")
    userID := session.Get("userID")

    if userID == nil {
        c.Redirect(http.StatusFound, "/login")
        return
    }

    id, ok := userID.(int)
    if !ok {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "code":    500,
            "message": "ID pengguna tidak valid",
            "title":   "Profile Akses Ditolak",
        })
        return
    }

    // Bind form checkout (metode pembayaran, alamat, dll)
    var input order.Peng
    if err := c.ShouldBind(&input); err != nil {
        c.HTML(http.StatusBadRequest, "keranjang.html", gin.H{
            "error":      "Semua kolom harus diisi",
            "title":      "Lengkapi Form",
            "isLoggedIn": true,
        })
        return
    }

    // Pastikan cart ada
    if cartData == nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "code":    500,
            "message": "Keranjang kosong",
            "title":   "Checkout",
        })
        return
    }

	cartJSON, ok := cartData.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format cart tidak valid"})
		return
	}

    // Decode JSON cart
    var cartItems []order.CartItem
    if err := json.Unmarshal([]byte(cartJSON), &cartItems); err != nil {
		log.Println("ERROR decode cart:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "Data keranjang tidak valid",
			"title":   "Checkout",
		})
		return
	}


	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Profile 1",
		})
		return
	}

	if user.Addres.ID == 0 {
		session.Set("success", "Alamat harus di isi")
		session.Save()
		c.Redirect(http.StatusFound, "/profile")
	}

    // Proses order ke database
    _ , err = h.orderService.CreateOrder(id, input, cartItems)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "code":    500,
            "message": "Gagal membuat order",
            "title":   "Checkout",
        })
        return
    }

    // Kosongkan keranjang
    session.Delete("cart")
    session.Save()

	session.Set("success", "Order berhasil")
	session.Save()

	c.Redirect(http.StatusFound,"/produk/transaction")
}

func (h *userHandler) IncreaseCartItem(c *gin.Context) {
    session := sessions.Default(c)
    cart := session.Get("cart")

    if cart == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang kosong"})
        return
    }

    itemID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
        return
    }


	cartJSON, ok := cart.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format cart tidak valid"})
		return
	}

    var items []order.CartItem
	if err := json.Unmarshal([]byte(cartJSON), &items); err != nil {
		log.Println("ERROR decode cart:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "Data keranjang tidak valid",
			"title":   "Checkout",
		})
		return
	}

	
    for i := range items {
        if items[i].ProductID == itemID {
            items[i].Quantity += 1
            items[i].SubTotal = items[i].Quantity * items[i].Price

			updatedCart, _ := json.Marshal(items)
   			session.Set("cart", string(updatedCart))
            session.Save()
			
			totalQty, totalPrice := helper.TotalCart(items)
            c.JSON(http.StatusOK, gin.H{
                "message": "Jumlah item bertambah",	
                "item":    items[i],
				"itemQuantity": items[i].Quantity,
				"itemSubtotal": items[i].SubTotal,
				"cart_count":    totalQty,
				"totalCart":  totalPrice,
				"isLoggedIn":  true,
            })
            return
        }
    }

	

     c.JSON(http.StatusNotFound, gin.H{"error": "Item tidak ditemukan"})
}

func (h *userHandler) DecreaseCartItem(c *gin.Context) {
    session := sessions.Default(c)
    cartData := session.Get("cart")

    if cartData == nil {
        c.JSON(400, gin.H{"error": "Cart kosong"})
        return
    }

    // cart := cartData.([]order.CartItem)
	
    itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
        return
    }

	cartJSON, ok := cartData.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format cart tidak valid"})
		return
	}

	var cart []order.CartItem
	if err := json.Unmarshal([]byte(cartJSON), &cart); err != nil {
		log.Println("ERROR decode cart:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "Data keranjang tidak valid",
			"title":   "Checkout",
		})
		return
	}

    var newItems []order.CartItem
    var updatedItem order.CartItem
    removed := false

    for _, item := range cart {
        if item.ProductID == itemID {

            item.Quantity--

            if item.Quantity <= 0 {
                removed = true
                continue // jangan masukkan ke newItems
            }

            updatedItem = item
            newItems = append(newItems, item)
        } else {
            newItems = append(newItems, item)
        }
    }

    // Simpan perubahan
	updatedCart, _ := json.Marshal(newItems)
   	session.Set("cart", string(updatedCart))
    session.Save()

    // Hitung total cart
    totalQty, totalPrice := helper.TotalCart(newItems)

    c.JSON(http.StatusOK, gin.H{
        "message":      "Item dikurangi",
        "removed":      removed,
        "itemQuantity": updatedItem.Quantity,
        "itemSubtotal": updatedItem.Price * int(updatedItem.Quantity),
        "cart_count":     totalQty,
        "totalCart":   totalPrice,
		"isLoggedIn":  true,
    })
}

func (h *userHandler) ViewTransaction(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")
	successMsg := session.Get("success")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Profile 1",
		})
		return
	}

	avatar := user.Avatar
	if  avatar ==""{
		avatar = "images/avatar/avatar.jpg"
	}

	var order []order.Order
	status := c.Query("status")
	if status != ""{
		order, err = h.orderService.GetByUserIdAndStatus(id,status)
		
	}else {
		order, err = h.orderService.GetOrderByUserId(id)
		
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
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

	c.HTML(http.StatusOK,"transaction.html", gin.H{
		"user":user,
		"avatar": avatar,
		"order":order,
		"statuskey":status,
		"isLoggedIn": true,
		"success":successMsg,
	})


}

func (h *userHandler) DetailTransaction(c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Profile 1",
		})
		return
	}

	avatar := user.Avatar
	if  avatar ==""{
		avatar = "images/avatar/avatar.jpg"
	}

	var uri order.UriInputOrderID
	err = c.ShouldBindUri(&uri)
	if err !=nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID order tidak valid",
			"title":   "input param id order error",
		})
		return
	}

	getDetail , err := h.orderService.GetDetailOrder(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat detail order",
			"title":   "get produk error",
		})
		return
	}

	order , err := h.orderService.GetOrderById(uri)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat detail order",
			"title":   "get produk error",
		})
		return
	}

	c.HTML(http.StatusOK,"detailOrder.html", gin.H{
		"user":user,
		"avatar": avatar,
		"isLoggedIn": true,
		"Detail":getDetail,
		"total" : order.TotPriceIDR(),
		"tanggal": order.CreatedAt,
	})
}

func (h *userHandler) Carapay( c *gin.Context){
	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    500,
			"message": "ID pengguna tidak valid",
			"title":   "Profile Akses Ditolak",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError,"error.html", gin.H{
			"code":    500,
			"message":"Gagal memuat user",
			"title":   "Profile 1",
		})
		return
	}

	avatar := user.Avatar
	if  avatar ==""{
		avatar = "/images/avatar/avatar.jpg"
	}

	c.HTML(http.StatusOK,"caraPay.html", gin.H{
		"user":user,
		"avatar": avatar,
		"isLoggedIn": true,
	})
}