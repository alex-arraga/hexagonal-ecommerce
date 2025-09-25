package domain

import "github.com/google/uuid"

type CartItem struct {
	ProductID uuid.UUID
	Quantity  int64
}

type Cart struct {
	UserID uuid.UUID
	Items  []CartItem
}

// NewCart creates a new cart for a user
func NewCart(userID uuid.UUID) *Cart {
	return &Cart{
		UserID: userID,
		Items:  []CartItem{},
	}
}

func (c *Cart) AddItem(productId uuid.UUID, quantity int64) {
	// Check if the item already exists in the cart
	for i, item := range c.Items {
		if item.ProductID == productId {
			// If it exists, update the quantity
			c.Items[i].Quantity += quantity
			return
		}

		c.Items = append(c.Items, CartItem{
			ProductID: productId,
			Quantity:  quantity,
		})
	}
}

// RemoveItem removes an item from the cart by product ID
func (c *Cart) RemoveItem(productID uuid.UUID) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			// Remove the item from the slice
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = []CartItem{}
}
