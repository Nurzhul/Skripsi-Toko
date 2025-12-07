package handler

import (
	"net/http"
	"toko/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type sessionUserHandler struct {
	userSevice user.Service
}

func NewSessionHendler(userService user.Service) *sessionUserHandler{
	return &sessionUserHandler{userService}
}

func (h *sessionUserHandler) New(c *gin.Context){
	errorMsg := ""
	if c.Query("error") == "true" {
		errorMsg = "Email atau password salah"
	}

	registerSuccess := c.Query("register_success")
	var message string
	if registerSuccess == "true" {
		message = "Akun berhasil dibuat. Silakan login."
	}

	
	c.HTML(http.StatusOK, "login.html", gin.H{
		"error": errorMsg,
		"success": message,
	})
	
}

func (h *sessionUserHandler) CreateUser(c *gin.Context){
	var input user.InputLogin

	
	err := c.ShouldBind(&input)
	if err !=nil {
		c.Redirect(http.StatusFound,"/login")
		return
	}
	
	user, err := h.userSevice.Login(input)
	if err != nil || user.Role !="user" {
		c.Redirect(http.StatusFound, "/login?error=true")
		return
	}
	
	session := sessions.Default(c)
	session.Set("userID",user.ID)
	session.Save()

	c.Redirect(http.StatusFound,"/dasboard")
}

func (h *sessionUserHandler) NewAdmin(c *gin.Context){
	errorMsg := ""
	if c.Query("error") == "true" {
		errorMsg = "Email atau password salah"
	}

	registerSuccess := c.Query("register_success")
	var message string
	if registerSuccess == "true" {
		message = "Akun berhasil dibuat. Silakan login."
	}

	
	c.HTML(http.StatusOK, "login_admin.html", gin.H{
		"error": errorMsg,
		"success": message,
	})
	
}
func (h *sessionUserHandler) CreateAdmin(c *gin.Context){
	var input user.InputLogin

	
	err := c.ShouldBind(&input)
	if err !=nil {
		c.Redirect(http.StatusFound,"/loginAdmin")
		return
	}
	
	user, err := h.userSevice.Login(input)
	if err != nil || user.Role !="admin" {
		c.Redirect(http.StatusFound, "/loginAdmin?error=true")
		return
	}
	
	session := sessions.Default(c)
	session.Set("userID",user.ID)
	session.Set("userRole",user.Role)
	session.Save()

	c.Redirect(http.StatusFound,"/dasboardAdmin")
}

func (h *sessionUserHandler) Destroy(c *gin.Context){
	seesion := sessions.Default(c)
	seesion.Clear()
	seesion.Save()

	c.Redirect(http.StatusFound,"/login")
}