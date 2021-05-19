package admin

//ProductDetails :
type ProductDetails struct {
	ProductID          int64                 `json:"productID"`
	ProductName        string                `json:"productName"`
	ProductPrize       float64               `json:"productPrize"`
	ProductQuantity    int64                 `json:"productQuantity"`
	ProductDiscount    float64               `json:"productDiscount"`
	ProductDescription string                `json:"productDescription"`
	ProductImages      []ProductImageDetails `json:"productImage,omitempty"`
}

//ProductImageDetails :
type ProductImageDetails struct {
	Name         string `json:"imageName,omitempty"`
	Base64String string `json:"base64String,omitempty"`
}

//OrderDetails :
type OrderDetails struct {
	OrderID         int64   `json:"orderID,omitempty"`
	ShippingAddress string  `json:"shippingAddress,omitempty"`
	Phone           int64   `json:"phone,omitempty"`
	TotalPayment    float64 `json:"totalPayment,omitempty"`
	OrderDate       string  `json:"orderDate,omitempty"`
	OrderStatus     string  `json:"orderStatus,omitempty"`
}

//OrdersDetailList :
type OrdersDetailList struct {
	TotalOrders int64          `json:"totalOrders"`
	Orders      []OrderDetails `json:"orders,omitempty"`
}
