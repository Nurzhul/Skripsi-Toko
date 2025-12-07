package order

import (
	"errors"
	"time"
)



type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository,}
}

type Service interface {
	CreateOrder(userID int, pay Peng, cartItems []CartItem) (Order, error)
	GetDetailOrder(uri UriInputOrderID) ([]DetailOrder, error)
	GetAllOrder() ([]Order, error)
	GetOrderById(uri UriInputOrderID) (Order, error)
	GetOrderByUserId(userID int) ([]Order, error)
	UpdateStatusOrder(uri UriInputOrderID, input string) error
	UpdateStatusPayOrder(uri UriInputOrderID, input string) error
	GetByStatus(status string) ([]Order, error)
	GetByUserIdAndStatus(userID int, status string) ([]Order, error)

}
func (s *service) CreateOrder(userID int, pay Peng, cartItems []CartItem) (Order, error) {
	if len(cartItems) == 0 {
		return Order{}, errors.New("Keranjang Kosong")
	}

	var total int 
	var orderItems []DetailOrder

	for _,item := range cartItems{
		subtotal := int(item.Quantity) * item.Price
		total += subtotal

		orderItems = append(orderItems, DetailOrder{
			ProdukID: item.ProductID,
			ProdukName: item.ProductName,
			Price: item.Price,
			Quantity: item.Quantity,
			Image: item.ProdukImage,
			SubTotal: subtotal,
		})

	}
	order := Order{
		UserID: userID,
		TotalPrice: total,
		Status: "Diproses",
		StatusPayment: "Menunggu",
		Pengambilan: pay.Pengambilan,
		CreatedAt: time.Now(),
	}

	err := s.repository.Save(order, orderItems)
	if err != nil {
		return Order{}, err
	}

	return order, nil

}

func (s *service) GetDetailOrder(uri UriInputOrderID) ([]DetailOrder, error){
	orders, err := s.repository.FindDetailOrderByOrdeID(uri.ID)
	if err != nil {
		return []DetailOrder{}, err
	}

	return orders, nil
}

func (s *service) GetAllOrder() ([]Order, error){
	ordes, err := s.repository.FindAllOrder()
	if err != nil {
		return []Order{}, err
	}

	return ordes, nil
}

func (s *service) GetOrderById(uri UriInputOrderID) (Order, error){
	order , err := s.repository.FindByOrderID(uri.ID)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (s *service) GetOrderByUserId(userID int) ([]Order, error){
	orders, err := s.repository.FindByUserID(userID)
	if err != nil {
		return orders, err
	}

	return orders, nil
}

func (s *service) UpdateStatusOrder(uri UriInputOrderID, input string) error{
	err := s.repository.UpdateStatus(uri.ID,input)
	if err != nil {
		return err
	}
	return nil
}

func(s *service) UpdateStatusPayOrder(uri UriInputOrderID, input string) error{
	err := s.repository.UpdateStatusPay(uri.ID,input)
	if err!= nil {
		return err
	}

	return nil
}

func(s *service) GetByStatus(status string) ([]Order, error){
	order, err := s.repository.FindByStatus(status)
	if err != nil {
		return []Order{}, err
	}

	return order, nil
}

func(s *service) GetByUserIdAndStatus(userID int, status string) ([]Order, error){
	order , err := s.repository.FindByUserIdBystatus(userID, status)
	if err != nil {
		return []Order{}, err
	}

	return order, nil
}