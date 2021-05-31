package front

import (
	"database/sql"
	"fmt"
	"myshop-api/api/common"
	"myshop-api/api/data"
	"strconv"
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
func GetProductsDetail(productID, productIDs, lowStockOrder, newStock, discountFilter string, limit, minPrize, maxPrize int) (objProductsDetails []ProductDetails, err error) {
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

	if minPrize > 0 && maxPrize > 0 && maxPrize > minPrize {
		productFinalPrize := "(pd.product_prize * (100-pd.product_discount)/100)"
		whrStr = whrStr + fmt.Sprintf(" AND "+productFinalPrize+" >= %v AND "+productFinalPrize+" <= %v ", minPrize, maxPrize)
	}

	discountParameters := strings.Split(discountFilter, ":")
	if len(discountParameters) == 3 {
		discount, _ := strconv.Atoi(discountParameters[1])
		if discountParameters[2] == "more" && discount > 0 {
			whrStr = whrStr + fmt.Sprintf(" AND pd.product_discount >= %v ", discount)
		}
	}

	if newStock != "" {
		orderbyStr = " , pd.dateadded DESC "
	}

	if whrStr != "" {
		whrStr = strings.Replace(whrStr, "AND", "WHERE", 1)
	}

	if orderbyStr != "" {
		orderbyStr = strings.Replace(orderbyStr, ",", "ORDER BY", 1)
	}

	if limit > 0 {
		limitStr = fmt.Sprintf(" LIMIT %v ", limit)
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

func RegisterNewUser(userName, userEmail, password string) error {
	sqlStr := fmt.Sprintf("INSERT INTO user (email_id, password, user_name) VALUES ('%v','%v','%v');", userEmail, password, userName)
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

//GetOrdersDetail :
func GetOrdersDetail(userID, orderType, searchOrderID string, limit, offset int) (objOrdersDetails []OrderDetails, err error) {
	var whrStr string
	if orderType != "" {
		whrStr = whrStr + " AND os.order_status_id = '" + orderType + "' "
	}
	if searchOrderID != "" {
		whrStr = whrStr + " AND order_id LIKE '%" + searchOrderID + "%' "
	}
	if userID != "" {
		whrStr = whrStr + " AND user_id = '" + userID + "' "
	}
	if whrStr != "" {
		whrStr = strings.Replace(whrStr, "AND", "WHERE", 1)
	}
	sqlStr := "SELECT order_id, shipping_address, phone, total_payment, order_date, os.status_name " +
		" FROM orders o " +
		" LEFT JOIN order_status os ON o.order_status = os.order_status_id " +
		whrStr + " LIMIT ?,?"

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

func GetTotalOrdersCount(userID, orderType, searchOrderID string) (totalRecords int64, err error) {

	var totolCount sql.NullInt64
	var whrStr string
	if orderType != "" {
		whrStr = whrStr + " AND order_status = '" + orderType + "' "
	}
	if searchOrderID != "" {
		whrStr = whrStr + " AND order_id LIKE '%" + searchOrderID + "%' "
	}
	if userID != "" {
		whrStr = whrStr + " AND user_id = '" + userID + "' "
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

//GetOrderFullDetail :
func GetOrderFullDetail(OrderID string) (objOrdersDetails OrderDetails, err error) {
	var whrStr string
	whrStr = " WHERE o.order_id = '" + OrderID + "' "

	sqlStr := "SELECT o.shipping_address, o.phone, o.total_payment, o.total_quantity, o.order_date, os.status_name, os.order_status_id, " +
		" p.product_name,od.quantity,od.prize,od.discount,u.email_id,u.user_name,GROUP_CONCAT(pim.image_name) " +
		" FROM orders o " +
		" LEFT JOIN order_detail od ON od.order_id = o.order_id " +
		" LEFT JOIN product_detail p ON od.product_id = p.product_id " +
		" LEFT JOIN `user` u ON u.user_id = o.user_id " +
		" LEFT JOIN order_status os ON o.order_status = os.order_status_id " +
		" LEFT JOIN product_images pim ON p.product_id = pim.product_id " + whrStr + "GROUP BY p.product_id ORDER BY o.order_date DESC "

	allRows, err := data.DemoDB.Query(sqlStr)
	if err != nil {
		return objOrdersDetails, err
	}
	for allRows.Next() {
		var objItemDetails OrderProductDetails
		var phoneNumber, totalQuantity, quantity, statusID sql.NullInt64
		var totalPayment, prize, discount sql.NullFloat64
		var orderDate, orderStatus, shippingAddress, productName, emailAddress, userName, productImage sql.NullString
		allRows.Scan(
			&shippingAddress,
			&phoneNumber,
			&totalPayment,
			&totalQuantity,
			&orderDate,
			&orderStatus,
			&statusID,
			&productName,
			&quantity,
			&prize,
			&discount,
			&emailAddress,
			&userName,
			&productImage,
		)

		objOrdersDetails.Username = userName.String
		objOrdersDetails.EmailID = emailAddress.String
		objOrdersDetails.ShippingAddress = shippingAddress.String
		objOrdersDetails.Phone = phoneNumber.Int64
		objOrdersDetails.TotalPayment = totalPayment.Float64
		objOrdersDetails.TotalQuantity = totalQuantity.Int64
		objOrdersDetails.OrderDate = orderDate.String
		objOrdersDetails.OrderStatus = orderStatus.String
		objOrdersDetails.OrderStatusID = statusID.Int64

		objItemDetails.ProductName = productName.String
		objItemDetails.Quantity = quantity.Int64
		objItemDetails.Prize = prize.Float64
		objItemDetails.Discount = discount.Float64

		productImages := strings.Split(productImage.String, ",")
		allProductImage := []ProductImageDetails{}
		for _, imageName := range productImages {
			var objProductImageDetails ProductImageDetails
			objProductImageDetails.Name = common.PRODUCT_IMAGE_PATH + imageName
			allProductImage = append(allProductImage, objProductImageDetails)
		}
		objItemDetails.ProductImages = allProductImage

		objOrdersDetails.OrderProductDetails = append(objOrdersDetails.OrderProductDetails, objItemDetails)
	}
	return objOrdersDetails, nil
}

func UpdateUserInfo(userID, userName string) error {
	sqlStr := fmt.Sprintf("UPDATE `user` u SET u.user_name = '%v' WHERE u.user_id = '%v';", userName, userID)

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

func UpdateUserPassword(userID, userPassword string) error {
	sqlStr := fmt.Sprintf("UPDATE `user` u SET u.password = %v WHERE u.user_id = %v;", userPassword, userID)

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
