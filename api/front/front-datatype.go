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
