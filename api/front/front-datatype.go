package front

//UserDetails :
type UserDetails struct {
	Username string `json:"username"`
	UserID   string `json:"userID"`
}

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

//UserOrderInput :
type UserOrderInput struct {
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	MobileNumber  int64           `json:"mobileNumber"`
	Address       string          `json:"address"`
	UserOrderInfo []UserOrderInfo `json:"userOrderInfo"`
}

//UserOrderInfo :
type UserOrderInfo struct {
	ProductID       int64 `json:"productID"`
	ProductQuantity int64 `json:"productQuantity"`
}

//OrderDetails :
type OrderDetails struct {
	OrderID             int64                 `json:"orderID,omitempty"`
	Username            string                `json:"username,omitempty"`
	EmailID             string                `json:"emailID,omitempty"`
	ShippingAddress     string                `json:"shippingAddress,omitempty"`
	Phone               int64                 `json:"phone,omitempty"`
	TotalPayment        float64               `json:"totalPayment"`
	TotalQuantity       int64                 `json:"totalQuantity"`
	OrderDate           string                `json:"orderDate,omitempty"`
	OrderStatus         string                `json:"orderStatus,omitempty"`
	OrderStatusID       int64                 `json:"orderStatusID,omitempty"`
	OrderProductDetails []OrderProductDetails `json:"productList,omitempty"`
}

//OrderProductDetails :
type OrderProductDetails struct {
	ProductName   string                `json:"productName,omitempty"`
	Quantity      int64                 `json:"quantity,omitempty"`
	Prize         float64               `json:"prize,omitempty"`
	Discount      float64               `json:"discount,omitempty"`
	ProductImages []ProductImageDetails `json:"productImage,omitempty"`
}

//OrdersDetailList :
type OrdersDetailList struct {
	TotalOrders int64          `json:"totalOrders"`
	Orders      []OrderDetails `json:"orders,omitempty"`
}
