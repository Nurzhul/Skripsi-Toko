package produk

import "gorm.io/gorm"

type Repository interface {
	Save(produk Produk) (Produk, error)
	FindByID(ID int) (Produk, error)
	GetAll() ([]Produk, error)
	Update(produk Produk)(Produk, error)
	Delete(ID int) error
	FindByCategory(category string)([]Produk,error)
	FindByName(name string)([]Produk, error)
	LastIDProduk() (Produk, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func (r *repository) Save(produk Produk) (Produk,error){
	err := r.db.Create(&produk).Error
	if err != nil{
		return produk, err
	}

	return produk, nil
}

func (r *repository) FindByID(ID int) (Produk, error){
	var produk Produk

	err := r.db.First(&produk, ID).Error
	if err != nil {
		return produk, err
	}

	return  produk, nil 
}

func (r *repository) GetAll()([]Produk, error){
	var produks []Produk

	err := r.db.Find(&produks).Error
	if err !=nil {
		return produks, err
	}

	return  produks, nil 
}

func (r *repository) Update(produk Produk)(Produk, error){
	err := r.db.Save(&produk).Error
	if err !=nil {
		return Produk{}, err
	}

	return produk, nil
}

func (r *repository) Delete(ID int) error{
	var produk Produk
	err := r.db.Delete(&produk, ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByCategory(category string)([]Produk, error){
	var produk []Produk
	err := r.db.Where("category LIKE ?","%"+category+"%").Find(&produk).Error
	if err != nil {
		return nil, err
	}

	return produk, nil
}

func (r *repository) FindByName(name string)([]Produk, error){
	var produk []Produk
	err := r.db.Where("name LIKE ?","%"+name+"%").Find(&produk).Error
	if err != nil {
		return nil, err
	}
	return produk, nil
}

func (r *repository) LastIDProduk()(Produk, error){
	var produk Produk
	err := r.db.Last(&produk).Error
	if err != nil {
		return Produk{}, err
	}

	return produk, nil
}
