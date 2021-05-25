package front

import (
	"database/sql"
	"fmt"
	"myshop-api/api/common"
	"myshop-api/api/data"
	"strings"
)

func CheckUserLoginDetails(userID, password string) (string, error) {
	var userName sql.NullString
	sqlStr := "SELECT user_name FROM `user` WHERE email_id = ? AND password = ?"
	err := data.DemoDB.QueryRow(sqlStr, userID, password).Scan(&userName)
	if err != nil && err != sql.ErrNoRows {
		return userName.String, err
	}
	return userName.String, nil
}

//GetProductsDetail :
func GetProductsDetail(productID, lowStockOrder, newStock string, limit int) (objProductsDetails []ProductDetails, err error) {
	var whrStr, orderbyStr, limitStr string
	if productID != "" {
		whrStr = whrStr + " AND pd.product_id = " + productID + " "
	}
	if lowStockOrder != "" {
		whrStr = whrStr + " AND pd.product_quantity <= 5 "
		orderbyStr = " , pd.product_quantity ASC "
	}

	if newStock != "" {
		orderbyStr = " , pd.dateadded DESC "
	}

	if limit > 0 {
		limitStr = fmt.Sprintf(" LIMIT %v ", limit)
	}

	if whrStr != "" {
		whrStr = strings.Replace(whrStr, "AND", "WHERE", 1)
	}

	if orderbyStr != "" {
		orderbyStr = strings.Replace(orderbyStr, ",", "ORDER BY", 1)
	}
	sqlStr := `SELECT pd.product_id,pd.product_name,pd.product_desc,pd.product_prize,pd.product_discount,pd.product_quantity, GROUP_CONCAT(pim.image_name) as product_images FROM product_detail pd 
	LEFT JOIN product_images pim ON pd.product_id = pim.product_id ` + whrStr + ` GROUP BY pd.product_id ` + orderbyStr + limitStr

	allRows, err := data.DemoDB.Query(sqlStr)
	if err != nil {
		return objProductsDetails, err
	}
	for allRows.Next() {
		var objProductDetails ProductDetails
		var allProductImage []ProductImageDetails
		var productID, productQuantity sql.NullInt64
		var productPrize, productDiscount sql.NullFloat64
		var productName, productDescription, productImage sql.NullString
		allRows.Scan(
			&productID,
			&productName,
			&productDescription,
			&productPrize,
			&productDiscount,
			&productQuantity,
			&productImage,
		)
		objProductDetails.ProductID = productID.Int64
		objProductDetails.ProductName = productName.String
		objProductDetails.ProductDescription = productDescription.String
		objProductDetails.ProductPrize = productPrize.Float64
		objProductDetails.ProductDiscount = productDiscount.Float64
		objProductDetails.ProductQuantity = productQuantity.Int64
		productImages := strings.Split(productImage.String, ",")

		for _, imageName := range productImages {
			var objProductImageDetails ProductImageDetails
			objProductImageDetails.Name = common.PRODUCT_IMAGE_PATH + imageName
			allProductImage = append(allProductImage, objProductImageDetails)
		}
		objProductDetails.ProductImages = allProductImage
		objProductsDetails = append(objProductsDetails, objProductDetails)
	}
	return objProductsDetails, nil
}
