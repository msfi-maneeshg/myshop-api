package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"myshop-api/api/common"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//AddProduct :
func AddProduct(w http.ResponseWriter, r *http.Request) {
	var objNewProductDetails ProductDetails
	var err error

	common.APIResponse(w, http.StatusOK, "Product has been added.")
	return
	//------check body request
	if r.Body == nil {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&objNewProductDetails)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Error:"+err.Error())
		return
	}

	//-------validate data
	if objNewProductDetails.ProductName == "" {
		common.APIResponse(w, http.StatusBadRequest, "Product Name can not be empty")
		return
	}

	if objNewProductDetails.ProductDescription == "" {
		common.APIResponse(w, http.StatusBadRequest, "Product Description can not be empty")
		return
	}

	//-------
	for index, productImage := range objNewProductDetails.ProductImages {
		fileName := time.Now().Format("20060102150405000") + "-" + fmt.Sprintf("%v", index) + productImage.Name + ".jpg"
		dec, err := base64.StdEncoding.DecodeString(productImage.Base64String)
		if err != nil {
			common.APIResponse(w, http.StatusBadRequest, "Invalid data of image."+err.Error())
			return
		}

		f, err := os.Create("assets/images/products/" + fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}

		productImage.Name = fileName
		objNewProductDetails.ProductImages[index] = productImage
	}

	objNewProductDetails.ProductID, err = InsertNewProductDetails(objNewProductDetails)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while inserting new product details. Error:"+err.Error())
		return
	}

	err = InsertNewProductImages(objNewProductDetails)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while inserting new product images. Error:"+err.Error())
		return
	}

	common.APIResponse(w, http.StatusOK, "Product has been added.")
}

//GetProductList:
func GetProductList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var productID = vars["productID"]

	productList, err := GetProductsDetail(productID)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while getting product list. Error:"+err.Error())
		return
	}

	if productID != "" {
		if len(productList) == 0 {
			common.APIResponse(w, http.StatusNotFound, "No product found!")
			return
		}
		common.APIResponse(w, http.StatusOK, productList[0])
		return
	}

	common.APIResponse(w, http.StatusOK, productList)
}

//EditProduct :
func EditProduct(w http.ResponseWriter, r *http.Request) {
	var objUpdatedProductDetails ProductDetails
	var err error

	vars := mux.Vars(r)
	var productID = vars["productID"]

	//------check body request
	if r.Body == nil {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&objUpdatedProductDetails)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Error:"+err.Error())
		return
	}

	//-------validate data
	if objUpdatedProductDetails.ProductName == "" {
		common.APIResponse(w, http.StatusBadRequest, "Product Name can not be empty")
		return
	}

	if objUpdatedProductDetails.ProductDescription == "" {
		common.APIResponse(w, http.StatusBadRequest, "Product Description can not be empty")
		return
	}

	//----validate productID
	productList, err := GetProductsDetail(productID)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while getting product list. Error:"+err.Error())
		return
	}
	if len(productList) == 0 {
		common.APIResponse(w, http.StatusNotFound, "No product found!")
		return
	}
	productIDInt, _ := strconv.Atoi(productID)
	objUpdatedProductDetails.ProductID = int64(productIDInt)
	//-------
	for index, productImage := range objUpdatedProductDetails.ProductImages {
		fileName := time.Now().Format("20060102150405000") + "-" + fmt.Sprintf("%v", index) + productImage.Name + ".jpg"
		dec, err := base64.StdEncoding.DecodeString(productImage.Base64String)
		if err != nil {
			common.APIResponse(w, http.StatusBadRequest, "Invalid data of image."+err.Error())
			return
		}

		f, err := os.Create("assets/images/products/" + fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}

		productImage.Name = fileName
		objUpdatedProductDetails.ProductImages[index] = productImage
	}

	err = UpdateProductDetails(objUpdatedProductDetails)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while updating product details. Error:"+err.Error())
		return
	}

	err = InsertNewProductImages(objUpdatedProductDetails)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while inserting product's new images. Error:"+err.Error())
		return
	}

	common.APIResponse(w, http.StatusOK, "Product details has been updated successfully.")
}
