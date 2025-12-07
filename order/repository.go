package order

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(order Order, items []DetailOrder) error
	FindByUserID(UserID int) ([]Order, error)
    FindByOrderID(OrderID int) (Order, error)
    FindDetailOrderByOrdeID(OrderId int) ([]DetailOrder, error)
    FindAllOrder()([]Order, error)
    Update(order Order) (Order, error)
    FindByStatus(status string) ([]Order, error)
    FindByUserIdBystatus(userID int, status string) ([]Order, error)
    UpdateStatus(id int, status string) error
    UpdateStatusPay(id int , statuspay string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func(r *repository) Save(order Order, items []DetailOrder) error{
	tx := r.db.Begin()

    if err := tx.Create(&order).Error; err != nil {
        tx.Rollback()
        return err
    }

    for _, item := range items {
        item.OrderID = order.ID
        if err := tx.Create(&item).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    return tx.Commit().Error
}

func(r *repository) FindByUserID(UserID int) ([]Order, error) {
	var orders []Order
	err := r.db.Where("user_id = ?", UserID).Order("id DESC").Find(&orders).Error
	if err != nil {
		return  orders, err
	}

	return  orders, nil 
}

func(r *repository) FindByOrderID(OrderID int) (Order, error){
    var order Order
    err := r.db.First(&order, OrderID).Error
    if err != nil {
        return order, err
    }

    return order, nil
}

func(r *repository)  FindDetailOrderByOrdeID(OrderId int) ([]DetailOrder, error){
    var detailorder []DetailOrder
    err := r.db.Where("order_id",OrderId).Find(&detailorder).Error
    if err != nil {
        return []DetailOrder{}, err
    }

    return  detailorder, nil
}

func(r *repository) FindAllOrder() ([]Order, error){
    var orders []Order
    err := r.db.Order("id DESC").Find(&orders).Error
    if err != nil {
        return orders, err
    }

    return orders, nil 
}

func(r *repository) Update(order Order) (Order, error){
    err := r.db.Save(&order).Error
    if err != nil {
        return order, err
    }

    return order, nil
}

func(r *repository) FindByStatus(status string) ([]Order, error){
    var order []Order
    err:= r.db.Where("status LIKE ?","%"+status+"%").Find(&order).Error
    if err != nil {
        return []Order{}, err
    }

    return order, nil
}

func(r *repository) FindByUserIdBystatus(userID int, status string) ([]Order, error){
    var order []Order
    err := r.db.Where("user_id", userID).Where("status LIKE ?","%"+status+"%").Find(&order).Error
    if err != nil {
        return []Order{}, err
    }

    return order, nil
}

func(r *repository) UpdateStatus(id int, status string) error{
    err := r.db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
    if err != nil {
        return err
    }

    return nil
}

func(r *repository) UpdateStatusPay(id int , statuspay string) error{
    err := r.db.Model(&Order{}).Where("id = ?",id).Update("status_payment",statuspay).Error
    if err != nil {
        return err
    }

    return nil
}