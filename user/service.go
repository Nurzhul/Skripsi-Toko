package user

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	RegisterUser(input InputRegister) (User, error)
	Login(input InputLogin) (User , error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
	GetAllUser()([]User, error)
	UpdateUser(input FromUpdateUser)(User,error)
	Updatepassword(ID int, input ChangePassword)(User, error)
	SaveAddres(ID int, input InputAddres)(Addres,error)
	UpdateAddres(ID int, input InputAddres)(Addres, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service{
	return &service{repository}
}

func (s *service) RegisterUser(input InputRegister) (User, error) {
	user := User{}
	emailcek, err := s.repository.FindByEmail(input.Email)
	if err == nil && emailcek.ID != 0 {
		return user, fmt.Errorf("email %s sudah digunakan ", input.Email)
	}else if err != nil && emailcek.ID == 0 {
		return User{}, err
	}
	
	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password),14)
	if err != nil {
		return User{}, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return User{}, err
	}

	return newUser, nil 
}

func (s *service) Login(input InputLogin)(User, error){
	email := input.Email
	password := input.Password

	user , err := s.repository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, errors.New("email tidak ditemukan")
		}
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(password))
	if err != nil {
		return User{} , errors.New("password salah ")
	}

	return user, nil 
}

func (s *service) SaveAvatar(ID int, fileLocation string )(User, error){
	user, err := s.repository.FindByID(ID)
	if err !=nil {
		return User{}, err
	}

	user.Avatar = fileLocation

	updateUser, err :=s.repository.Update(user)
	if err != nil {
		return User{}, err
	}

	return updateUser, nil
}

func (s *service) GetUserByID(ID int) (User,error){
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return User{},err
	}

	if user.ID == 0 {
		return User{}, errors.New("User Tidak di temukan ")
	}

	return user,nil
}

func (s *service) GetAllUser()([]User,error){
	users, err := s.repository.FindAll()
	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func (s *service) UpdateUser(input FromUpdateUser) (User, error){
	user , err := s.repository.FindByID(input.ID)
	if err != nil {
		return User{}, err
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone

	updateUser , err := s.repository.Update(user)
	if err != nil {
		return User{}, err
	}

	return updateUser, nil
}

func (s *service) Updatepassword(ID int , input ChangePassword)(User, error){
	user, err := s.repository.FindByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return User{}, errors.New("user tidak ditemukan")
		}
		return User{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(input.OldPassword));err!= nil{
		return User{},errors.New("password lama salah")
	} 	

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword),14)
	if err!= nil {
		return User{}, err
	}

	user.PasswordHash = string(hashedPassword)
	return s.repository.Update(user)
}

func (s *service)SaveAddres(ID int,input InputAddres)(Addres,error){
	user, err := s.repository.FindByID(ID)
	if err !=nil {
		return Addres{}, err
	}
	addres := Addres{}
	addres.Provinsi = input.Provinsi
	addres.Kabupaten = input.Kabupaten
	addres.Kecamatan = input.Kecamatan
	addres.Desa = input.Desa
	addres.Jalan = input.Jalan
	addres.Rt = input.Rt
	addres.Rw = input.Rw
	addres.NoRm = input.NoRm
	addres.KodePos = input.KodePos
	addres.Deskripsi = input.Deskripsi
	addres.UserID = user.ID

	newAddres , err := s.repository.SaveAddres(addres)
	if err != nil {
		return Addres{}, err
	}

	return newAddres, nil
}

func (s *service)UpdateAddres(ID int, input InputAddres)(Addres, error){
	user, err := s.repository.FindByID(ID)
	if err !=nil {
		return Addres{}, err
	}
	addres := Addres{}
	addres.Provinsi = input.Provinsi
	addres.Kabupaten = input.Kabupaten
	addres.Kecamatan = input.Kecamatan
	addres.Desa = input.Desa
	addres.Jalan = input.Jalan
	addres.Rt = input.Rt
	addres.Rw = input.Rw
	addres.NoRm = input.NoRm
	addres.KodePos = input.KodePos
	addres.Deskripsi = input.Deskripsi
	addres.UserID = user.ID
	
	newAddres , err := s.repository.UpdateAddres(addres)
	if err != nil {
		return Addres{}, err
	}

	return newAddres, nil
}