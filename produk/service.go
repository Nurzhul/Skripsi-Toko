package produk

type Service interface {
	AddProduk(input InputAddProduk, imageFile string) (Produk, error)
	GetProdukById(Uri GetProdukUri) (Produk, error)
	GetAllProduk() ([]Produk, error)
	UpdateProduk(ID GetProdukUri, input InputUpdateProduk, imageFile string) (Produk, error)
	DeleteProduk(ID GetProdukUri) error
	GetByCategory(category string) ([]Produk, error)
	GetByName(name string) ([]Produk, error)
	GetlastProdukID() (int, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) AddProduk(input InputAddProduk, imageFile string) (Produk, error) {
	produk := Produk{
		Name:        input.Name,
		Description: input.Description,
		Category:    input.Category,
		Stock:       input.Stock,
		Price:       input.Price,
		ImageName:   imageFile,
	}

	newProduk, err := s.repository.Save(produk)
	if err != nil {
		return Produk{}, err
	}

	return newProduk, nil
}

func (s *service) GetProdukById(Uri GetProdukUri) (Produk, error) {
	produk, err := s.repository.FindByID(Uri.ID)
	if err != nil {
		return Produk{}, err
	}

	return produk, nil
}

func (s *service) GetAllProduk() ([]Produk, error) {
	produks, err := s.repository.GetAll()
	if err != nil {
		return []Produk{}, err
	}

	return produks, nil
}

func (s *service) UpdateProduk(ID GetProdukUri, input InputUpdateProduk, imageFile string) (Produk, error) {
	produk, err := s.repository.FindByID(ID.ID)
	if err != nil {
		return Produk{}, err
	}

	if input.Name != "" {
		produk.Name = input.Name
	}
	if input.Description != "" {
		produk.Description = input.Description
	}
	if input.Category != "" {
		produk.Category = input.Category
	}
	if input.Stock != 0 {
		produk.Stock = input.Stock
	}
	if input.Price != 0 {
		produk.Price = input.Price
	}
	if imageFile != "" {
		produk.ImageName = imageFile
	}

	upProduk, err := s.repository.Update(produk)
	if err != nil {
		return Produk{}, err
	}

	return upProduk, nil
}

func (s *service) DeleteProduk(ID GetProdukUri) error {
	produk, err := s.repository.FindByID(ID.ID)
	if err != nil {
		return err
	}

	return s.repository.Delete(produk.ID)

}

func (s *service) GetByCategory(category string) ([]Produk, error) {
	produk, err := s.repository.FindByCategory(category)
	if err != nil {
		return nil, err
	}

	return produk, nil

}

func (s *service) GetByName(name string) ([]Produk, error) {
	produk, err := s.repository.FindByName(name)
	if err != nil {
		return nil, err
	}

	return produk, nil
}

func (s *service) GetlastProdukID() (int, error) {
	lastID, err := s.repository.LastIDProduk()
	if err != nil {
		return 0, err
	}

	return lastID.ID, nil

}