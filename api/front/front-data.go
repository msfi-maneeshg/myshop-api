package front

import (
	"database/sql"
	"fmt"
	"myshop-api/api/common"
	"myshop-api/api/data"
	"strings"
)

func CheckUserLoginDetails(userEmail, password string) (objUserInfo UserDetails, err error) {
	var userName, userID sql.NullString
	sqlStr := "SELECT user_name,user_id FROM `user` WHERE email_id = ?"
	if password != "" {
		sqlStr = "SELECT user_name,user_id FROM `user` WHERE email_id = ? AND password = '" + password + "'"
	}
	err = data.DemoDB.QueryRow(sqlStr, userEmail).Scan(&userName, &userID)
	if err != nil && err != sql.ErrNoRows {
		return objUserInfo, err
	}

	objUserInfo.Username = userName.String
	objUserInfo.UserID = userID.String
	return objUserInfo, nil
}

//GetProductsDetail :
func GetProductsDetail(productID, productIDs, lowStockOrder, newStock string, limit int) (objProductsDetails []ProductDetails, err error) {
	var whrStr, orderbyStr, limitStr string
	if productID != "" {
		whrStr = whrStr + " AND pd.product_id = " + productID + " "
	} else if productIDs != "" {
		whrStr = whrStr + " AND pd.product_id IN (" + productIDs + ") "
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

//InsertUserOrder :
func InsertUserOrder(objUserInfo UserDetails, objUserOrderInput UserOrderInput, totalPayment float64) (orderID int64, err error) {

	sqlStr := fmt.Sprintf("INSERT INTO orders (user_id, shipping_address, phone, total_payment, total_quantity) VALUES ('%v','%v','%v','%v','%v')", objUserInfo.UserID, objUserOrderInput.Address, objUserOrderInput.MobileNumber, totalPayment, len(objUserOrderInput.UserOrderInfo))

	stmt, err := data.DemoDB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return orderID, err
	}

	res, err := stmt.Exec()
	if err != nil {
		return orderID, err
	}
	orderID, err = res.LastInsertId()
	if err != nil {
		return orderID, err
	}
	return orderID, nil
}

//InsertUserOrderDetails :
func InsertUserOrderDetails(orderID int64, objUserInfo UserDetails, objProductDetails []ProductDetails, totalPayment float64) (err error) {
	sqlStr := "INSERT INTO order_detail (order_id, product_id, quantity, prize, discount) VALUES "
	for _, productInfo := range objProductDetails {
		sqlStr = sqlStr + fmt.Sprintf("('%v','%v','%v','%v','%v'),", orderID, productInfo.ProductID, productInfo.ProductQuantity, productInfo.ProductPrize, productInfo.ProductDiscount)
	}

	if sqlStr == "" {
		return nil
	}
	sqlStr = sqlStr[:len(sqlStr)-1]

	stmt, err := data.DemoDB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func UpdateProductStock(productInfo ProductDetails) (err error) {

	sqlStr := fmt.Sprintf("UPDATE product_detail p SET p.product_quantity = p.product_quantity - %v WHERE p.product_id = %v;", productInfo.ProductQuantity, productInfo.ProductID)

	if sqlStr == "" {
		return nil
	}

	stmt, err := data.DemoDB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
