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

func GetProductsDetail() (objProductsDetails []ProductDetails, err error) {
	sqlStr := `SELECT pd.product_id,pd.product_name,pd.product_desc,pd.product_prize,pd.product_discount,pd.product_quantity, GROUP_CONCAT(pim.image_name) as product_images FROM product_detail pd 
	LEFT JOIN product_images pim ON pd.product_id = pim.product_id 
	GROUP BY pd.product_id`

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
