package admin

import (
	"database/sql"
	"fmt"
	"myshop-api/api/common"
	"myshop-api/api/data"
	"strings"
)

//InsertNewProductDetails :
func InsertNewProductDetails(objNewProductDetails ProductDetails) (newProductID int64, err error) {
	sqlStr := fmt.Sprintf("INSERT INTO product_detail (product_name, product_desc, product_prize, product_discount, product_quantity) VALUES ('%v','%v','%v','%v','%v')", objNewProductDetails.ProductName, objNewProductDetails.ProductDescription, objNewProductDetails.ProductPrize, objNewProductDetails.ProductDiscount, objNewProductDetails.ProductQuantity)

	stmt, err := data.DemoDB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return newProductID, err
	}

	res, err := stmt.Exec()
	if err != nil {
		return newProductID, err
	}

	newProductID, err = res.LastInsertId()
	if err != nil {
		return newProductID, err
	}
	return newProductID, nil
}

//InsertNewProductImages :
func InsertNewProductImages(objNewProductDetails ProductDetails) (err error) {
	sqlStr := "INSERT INTO product_images (product_id, image_name) VALUES "
	var sqlSubStr string
	for _, imageInfo := range objNewProductDetails.ProductImages {
		sqlSubStr = sqlSubStr + fmt.Sprintf("('%v','%v'),", objNewProductDetails.ProductID, imageInfo.Name)
	}

	if sqlSubStr == "" {
		return nil
	}
	sqlSubStr = sqlSubStr[:len(sqlSubStr)-1]
	stmt, err := data.DemoDB.Prepare(sqlStr + sqlSubStr + ";")
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

//GetProductsDetail :
func GetProductsDetail(productID string) (objProductsDetails []ProductDetails, err error) {
	var whrStr string
	if productID != "" {
		whrStr = " WHERE pd.product_id = " + productID + " "
	}
	sqlStr := `SELECT pd.product_id,pd.product_name,pd.product_desc,pd.product_prize,pd.product_discount,pd.product_quantity, GROUP_CONCAT(pim.image_name) as product_images FROM product_detail pd 
	LEFT JOIN product_images pim ON pd.product_id = pim.product_id ` + whrStr + ` GROUP BY pd.product_id`

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

//UpdateProductDetails :
func UpdateProductDetails(objUpdatedProductDetails ProductDetails) error {
	sqlStr := fmt.Sprintf("UPDATE product_detail SET product_name = '%v', product_desc = '%v',product_prize = '%v',product_discount = '%v',product_quantity = '%v' where product_id = '%v'; ", objUpdatedProductDetails.ProductName, objUpdatedProductDetails.ProductDescription, objUpdatedProductDetails.ProductPrize, objUpdatedProductDetails.ProductDiscount, objUpdatedProductDetails.ProductQuantity, objUpdatedProductDetails.ProductID)

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

//GetOrdersDetail :
func GetOrdersDetail(orderType, searchOrderID string, limit, offset int) (objOrdersDetails []OrderDetails, err error) {
	var whrStr string
	if orderType != "" {
		whrStr = whrStr + " AND order_status = '" + orderType + "' "
	}
	if searchOrderID != "" {
		whrStr = whrStr + " AND order_id LIKE '%" + searchOrderID + "%' "
	}
	if whrStr != "" {
		whrStr = strings.Replace(whrStr, "AND", "WHERE", 1)
	}
	sqlStr := "SELECT order_id, shipping_address, phone, total_payment, order_date, order_status " +
		" FROM orders " + whrStr + " LIMIT ?,?"

	allRows, err := data.DemoDB.Query(sqlStr, offset, limit)
	if err != nil {
		return objOrdersDetails, err
	}
	for allRows.Next() {
		var objOrderDetails OrderDetails
		var orderID, phoneNumber sql.NullInt64
		var totalPayment sql.NullFloat64
		var orderDate, orderStatus, shippingAddress sql.NullString
		allRows.Scan(
			&orderID,
			&shippingAddress,
			&phoneNumber,
			&totalPayment,
			&orderDate,
			&orderStatus,
		)
		objOrderDetails.OrderID = orderID.Int64
		objOrderDetails.ShippingAddress = shippingAddress.String
		objOrderDetails.Phone = phoneNumber.Int64
		objOrderDetails.TotalPayment = totalPayment.Float64
		objOrderDetails.OrderDate = orderDate.String
		objOrderDetails.OrderStatus = orderStatus.String

		objOrdersDetails = append(objOrdersDetails, objOrderDetails)
	}
	return objOrdersDetails, nil
}

func GetTotalOrdersCount(orderType, searchOrderID string) (totalRecords int64, err error) {

	var totolCount sql.NullInt64
	var whrStr string
	if orderType != "" {
		whrStr = whrStr + " AND order_status = '" + orderType + "' "
	}
	if searchOrderID != "" {
		whrStr = whrStr + " AND order_id LIKE '%" + searchOrderID + "%' "
	}
	if whrStr != "" {
		whrStr = strings.Replace(whrStr, "AND", "WHERE", 1)
	}

	sqlStr := "SELECT COUNT(order_id) as total_records FROM orders"
	err = data.DemoDB.QueryRow(sqlStr + whrStr).Scan(&totolCount)
	if err != nil && err != sql.ErrNoRows {
		return totalRecords, err
	}
	totalRecords = totolCount.Int64
	return totalRecords, nil
}
