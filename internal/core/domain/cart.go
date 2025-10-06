package domain

import "github.com/google/uuid"

type CartItem struct {
	ProductID uuid.UUID
	Quantity  int16
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

func (c *Cart) AddItem(productId uuid.UUID, quantity int16) error {
	// Check if the item already exists in the cart
	for i, item := range c.Items {
		if item.ProductID == productId {
			// If it exists, update the quantity
			c.Items[i].Quantity += quantity

			// Validate if quantity is less than 0
			if c.Items[i].Quantity <= 0 {
				c.RemoveItem(productId)
				return nil
			}
			return nil
		}
	}

	// If product not exist and quantity is negative, return error
	if quantity < 0 {
		return ErrNegativeQuantityNonExistProductCart
	}

	c.Items = append(c.Items, CartItem{
		ProductID: productId,
		Quantity:  quantity,
	})
	return nil
}

// RemoveItem removes an item from the cart by product ID
func (c *Cart) RemoveItem(productID uuid.UUID) error {
	if len(c.Items) <= 0 {
		return ErrAlreadyEmptyCart
	}

	for i, item := range c.Items {
		// Check if productID exist in cart
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...) // Remove the item from the slice
			return nil
		}
	}
	
	return ErrProductNotFoundCart
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = []CartItem{}
}
