package helper

import "toko/order"

func TotalCart(cart []order.CartItem) (totalQty int, totalPrice int) {
	for _, item := range cart {
		totalQty += item.Quantity
		totalPrice += int(item.Quantity) * item.Price
	}
	return
}