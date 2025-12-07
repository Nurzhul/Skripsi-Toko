package user

type InputRegister struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Phone    string `form:"phone" binding:"required"`
	Password string `form:"password" binding:"required,min=8"`
}

type InputLogin struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type FromUpdateUser struct {
	ID    int
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required,email"`
	Phone string `form:"phone" binding:"required"`
}

type ChangePassword struct {
	OldPassword     string `form:"old_password" binding:"required"`
	NewPassword     string `form:"new_password" binding:"required,min=8"`
	CpmfirmPassword string `form:"comfirm_password" binding:"required"`
}

type InputAddres struct {
	Provinsi  string `form:"provinsi" binding:"required"`
	Kabupaten string `form:"kabupaten" binding:"required"`
	Kecamatan string `form:"kecamatan" binding:"required"`
	Desa      string `form:"desa" binding:"required"`
	Jalan     string `form:"jalan" binding:"required"`
	Rt        string `form:"rt" binding:"required"`
	Rw        string `form:"rw" binding:"required"`
	NoRm      string `form:"no_rm" binding:"required"`
	KodePos   string `form:"kode_pos" binding:"required"`
	Deskripsi string `form:"deskripsi" binding:"required"`
}