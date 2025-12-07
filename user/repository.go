package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string)(User, error)
	FindByID(ID int)(User, error)
	Update(user User)(User, error)
	FindAll()([]User, error)
	FindName(name string)(User , error)
	SaveAddres(addres Addres)(Addres,error)
	UpdateAddres(addres Addres)(Addres,error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func (r *repository) Save(user User) (User, error){
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error){
	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil 
}

func(r *repository) FindByID(ID int) (User, error){
	var user User

	err := r.db.Preload("Addres").First(&user, ID).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) Update(user User) (User, error){
	err := r.db.Save(&user).Error
	if err != nil {
		return user , err
	}

	return user, nil
}

func (r *repository) FindAll()([]User, error){
	var users []User

	err := r.db.Preload("Addres").Find(&users).Error

	if err != nil {
		return users, err
	}

	return users, nil
}

func (r *repository) FindName(name string) (User, error) {
	var user User
	err := r.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, err
}

func(r *repository) SaveAddres(addres Addres)(Addres,error){
	err := r.db.Create(&addres).Error
	if err != nil {
		return  addres, err
	}

	return addres,nil
}

func(r *repository) UpdateAddres(addres Addres)(Addres,error){
	err :=r.db.Save(&addres).Error
	if err!= nil {
		return addres, err
	}

	return addres, nil
}