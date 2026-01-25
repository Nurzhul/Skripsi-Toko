package main

import (
	"log"
	"net/http"
	"os"
	"toko/handler"
	"toko/order"
	"toko/produk"
	"toko/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}


	dsn :=os.Getenv("DNS")
	key :=os.Getenv("SECRET_KEY")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// err = db.AutoMigrate(
	// 	&userEntity.User{},
	// 	&userEntity.Addres{},
	// 	&produkEntity.Produk{},
	// 	&orderEntity.Order{},
	// 	&orderEntity.DetailOrder{},
	// )
	
	// if err != nil {
    //     log.Fatal("Gagal migrate: ", err)
    // }


	userRepository := user.NewRepository(db)
	produkRepository := produk.NewRepository(db)
	orderRepository := order.NewRepository(db)
	
	userService := user.NewService(userRepository)
	produkService := produk.NewService(produkRepository)
	orderService :=order.NewService(orderRepository)

	SeesionHandler := handler.NewSessionHendler(userService)
	userHandler := handler.NewUserHandler(userService,produkService,orderService)
	adminHandler := handler.NewAdminHandler(userService,produkService,orderService)
	

	// gin.SetMode(gin.ReleaseMode) // kalau untuk production -> perlucaritahu dulu 
	router := gin.Default()

	cookieStore := cookie.NewStore([]byte(key))
	cookieStore.Options(sessions.Options{
		Path: "/",
		MaxAge: 60 * 60 * 24 * 7, // 7 hari
		HttpOnly: true,
	})
	router.Use(sessions.Sessions("tokoPertanian", cookieStore))

	router.LoadHTMLGlob("template/**/*.html")
	router.Static("/images","./images")
	router.Static("/css","./asset/css")
	router.Static("/js","./asset/js")

	// user
	router.GET("/dasboard",userHandler.Index)
	router.GET("/register",userHandler.New)
	router.POST("/register",userHandler.Regis)
	router.GET("/profile",AuthUserMiddleware(),userHandler.Profile)
	router.GET("/profile/password",AuthUserMiddleware(),userHandler.Newpassword)
	router.POST("/profile/password",AuthUserMiddleware(),userHandler.UpdatePass)
	router.GET("/profile/addres",AuthUserMiddleware(),userHandler.Addres)
	router.POST("/profile/addres",AuthUserMiddleware(),userHandler.NewAddres)
	router.GET("/profile/saveAvatar", AuthUserMiddleware(),userHandler.SaveGetAvatar)
	router.POST("/profile/saveAvatar",AuthUserMiddleware(),userHandler.SavePostAvatar)
	router.GET("/produk",userHandler.Produk)
	router.GET("/produk/:id",userHandler.DetailProduk)
	router.POST("/produk/addCart",AuthUserMiddleware(),userHandler.AddToCart)
	router.GET("/produk/Keranjang",AuthUserMiddleware(),userHandler.ViewCart)
	router.GET("/produk/cartCount",AuthUserMiddleware(),userHandler.GetCartCount)
	router.POST("/cart/increase/:id",AuthUserMiddleware(),userHandler.IncreaseCartItem)
	router.POST("/cart/decrease/:id",AuthUserMiddleware(),userHandler.DecreaseCartItem)
	router.POST("/cart/remove/:id", AuthUserMiddleware(),userHandler.RemoveCartItem)
	// router.POST("/cart/decrease/:itemID",AuthUserMiddleware(),userHandler.DecreaseCartItem)
	router.POST("/cart/order",AuthUserMiddleware(),userHandler.Checkout)
	router.GET("/produk/transaction",AuthUserMiddleware(),userHandler.ViewTransaction)
	router.GET("/produk/transaction/:id",AuthUserMiddleware(),userHandler.DetailTransaction)
	router.GET("/produk/transaction/cara", AuthUserMiddleware(),userHandler.Carapay)
	// admin 
	router.GET("/dasboardAdmin",AuthAdminMiddleware(),adminHandler.IndexAdmin)
	router.GET("/dasboardAdmin/user/:id", AuthAdminMiddleware(),adminHandler.GetDetailUser)
	router.GET("/produkAdmin",AuthAdminMiddleware(),adminHandler.ShouProduk)
	router.GET("/produkAdmin/addProduk",AuthAdminMiddleware(),adminHandler.AddGetProduk)
	router.POST("/produkAdmin/addProduk",AuthAdminMiddleware(),adminHandler.AddPostProduk)
	router.GET("/produkAdmin/updateProduk/:id",AuthAdminMiddleware(),adminHandler.UpdateGetProduk)
	router.POST("/produkAdmin/updateProduk/:id",AuthAdminMiddleware(),adminHandler.UpdatePostProduk)
	router.POST("/produkAdmin/deleteProduk/:id",AuthAdminMiddleware(),adminHandler.DeletePostProduk)
	router.GET("/produkAdmin/detail/:id",AuthAdminMiddleware(),adminHandler.DetailProduk)
	router.GET("/produkAdmin/Transaksi",AuthAdminMiddleware(),adminHandler.GetLihatTransaksi)
	router.POST("/produkAdmin/Transaksi/updateStatus/:id",AuthAdminMiddleware(),adminHandler.PostUpdateStatus)
	router.POST("/produkAdmin/Transaksi/updateStatusPay/:id",AuthAdminMiddleware(),adminHandler.PostUpdateStatusPay)
	router.GET("/produkAdmin/Transaksi/:id",AuthAdminMiddleware(),adminHandler.GetDetailTransaction)

	//session
	router.GET("/login", SeesionHandler.New)
	router.POST("/session",SeesionHandler.CreateUser)
	router.GET("/logout",SeesionHandler.Destroy)
	router.GET("/loginAdmin",SeesionHandler.NewAdmin)
	router.POST("/sessionAdmin",SeesionHandler.CreateAdmin)

	
	router.Run()
}

func AuthUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID")

		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userID := session.Get("userID")
		userRole := session.Get("userRole")

		if userID == nil || userRole != "admin" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}