package controllers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"
	"myapp/models"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestHomeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/", HomeHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	mockResponse := "Welcome to the Product API"
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, resp.Code)
}
func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `products`").
		WithArgs("APPLE", 99.0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/products", CreateProduct)

	productJSON := `{"name": "APPLE", "price": 99.0}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(productJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	expectedBody := `{"message":"Product created successfully","id":1}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `products` WHERE `products`.`id` = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.DELETE("/products/:id", DeleteProduct)

	req, err := http.NewRequest("DELETE", "/products/1", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	expectedBody := `{"message":"Product deleted successfully"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	expectedQuery := "SELECT * FROM `products` WHERE `products`.`id` = ? ORDER BY `products`.`id` LIMIT ?"
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(1, "APPLE", 22.0))

	expectedUpdate := "UPDATE `products` SET `name`=?,`price`=? WHERE `id` = ?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedUpdate)).
		WithArgs("APPLE", 100.0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.PUT("/products/:id", UpdateProduct)

	productJSON := `{"name": "APPLE", "price": 100.0}`
	req, err := http.NewRequest("PUT", "/products/1", bytes.NewBufferString(productJSON))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	expectedBody := `{"message":"Product updated successfully"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestGetAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery(`^SELECT \* FROM ` + "`products`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(1, "APPLE", 99.0).
			AddRow(2, "BANANA", 50.0))

	models.DB = gormDB

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/products", GetAllProducts)

	req, _ := http.NewRequest("GET", "/products", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	expectedBody := `[{"id":1,"name":"APPLE","price":99.0},{"id":2,"name":"BANANA","price":50.00}]`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestGetProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	expectedQuery := "SELECT * FROM `products` WHERE `products`.`id` = ? ORDER BY `products`.`id` LIMIT ?"
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(1, "APPLE", 99.0))

	models.DB = gormDB

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/products/:id", GetProductByID)

	req, _ := http.NewRequest("GET", "/products/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	expectedBody := `{"product":{"id":1,"name":"APPLE","price":99}}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

// Test for abnormal process
func TestCreateProductWithInvalidData(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `products`").
		WithArgs("APPLE", 10.0).
		WillReturnError(errors.New("Validation error"))
	mock.ExpectRollback()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/products", CreateProduct)

	productJSON := `{"name": "APPLE", "price": "invalid_price"}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(productJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	expectedBody := `{"error":"Invalid input"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestCreateProductWithLackData(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `products`").
		WithArgs("", 10.0). // Invalid data
		WillReturnError(errors.New("Validation error"))
	mock.ExpectRollback()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/products", CreateProduct)

	productJSON := `{"name": "", "price": 10.0}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(productJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	expectedBody := `{"error":"Name and Price are required"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestCreateProductFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `products`").
		WithArgs("APPLE", 99.0).
		WillReturnError(errors.New("Failed to insert into database"))
	mock.ExpectRollback()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/products", CreateProduct)

	validProductJSON := `{"name": "APPLE", "price": 99.0}`
	req, err := http.NewRequest("POST", "/products", bytes.NewBufferString(validProductJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	expectedBody := `{"error":"Failed to create product"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestDeleteProductFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB
	// Mock Error Delete Product
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `products` WHERE `products`.`id` = ?")).
		WithArgs(999).
		WillReturnError(errors.New("Failed to delete product"))
	mock.ExpectRollback()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.DELETE("/products/:id", DeleteProduct)

	req, err := http.NewRequest("DELETE", "/products/99", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	expectedBody := `{"error":"Failed to delete product"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}
func TestUpdateProductFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `products` WHERE `products`.`id` = ? ORDER BY `products`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(1, "APPLE", 50.0))

	// Mock Error Update Product
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `products` SET `name`=?,`price`=? WHERE `id` = ?")).
		WithArgs("APPLE_UPDATED", 200.0, 1).
		WillReturnError(errors.New("update failed"))
	mock.ExpectRollback()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.PUT("/products/:id", UpdateProduct)

	productJSON := `{"name": "APPLE_UPDATED", "price": 200.0}`
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBufferString(productJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	expectedBody := `{"error":"Failed to update product"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}
func TestGetProductByIDWithInvalidProductID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	expectedQuery := "SELECT * FROM `products` WHERE `products`.`id` = ?"
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(1, 1).
		WillReturnError(errors.New("Invalid product ID"))

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/products/:id", GetProductByID)

	req, err := http.NewRequest("GET", "/products/Invalid", nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	expectedBody := `{"error":"Invalid product ID"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestGetProductByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	models.DB = gormDB

	expectedQuery := "SELECT * FROM `products` WHERE `products`.`id` = ? ORDER BY `products`.`id` LIMIT ?"
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(1, 1).
		WillReturnError(errors.New("Product not found"))

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/products/:id", GetProductByID)

	req, err := http.NewRequest("GET", "/products/1", nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	expectedBody := `{"error":"Product not found"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}
