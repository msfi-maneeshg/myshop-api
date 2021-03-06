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
	ProductName string  `json:"productName,omitempty"`
	Quantity    int64   `json:"quantity,omitempty"`
	Prize       float64 `json:"prize,omitempty"`
	Discount    float64 `json:"discount,omitempty"`
}

//OrdersDetailList :
type OrdersDetailList struct {
	TotalOrders int64          `json:"totalOrders"`
	Orders      []OrderDetails `json:"orders,omitempty"`
}

//UserDetails :
type UserDetails struct {
	Username string `json:"username"`
	UserID   string `json:"userID"`
}

//CategoryDetails :
type CategoryDetails struct {
	CategoryID   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	CategoryURL  string `json:"category_url"`
}

//CategoryDetailList :
type CategoryDetailList struct {
	TotalCategories int64             `json:"totalCategories"`
	Categories      []CategoryDetails `json:"categories,omitempty"`
}
