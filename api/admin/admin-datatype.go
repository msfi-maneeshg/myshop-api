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
