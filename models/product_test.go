package models_test

import (
	"errors"
	"myapp/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

	product := &models.Product{
		Name:  "APPLE",
		Price: 99.0,
	}

	id, err := models.CreateProduct(product)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
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

	rowsAffected, err := models.DeleteProduct(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowsAffected)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
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
			AddRow(1, "APPLE", 50.0))

	expectedUpdate := "UPDATE `products` SET `name`=?,`price`=? WHERE `id` = ?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedUpdate)).
		WithArgs("APPLE_UPDATED", 100.0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updatedProduct := &models.Product{
		Name:  "APPLE_UPDATED",
		Price: 100.0,
	}

	err = models.UpdateProduct(1, updatedProduct)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
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

	products, err := models.GetAllProducts()
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "APPLE", products[0].Name)
	assert.Equal(t, 99.0, products[0].Price)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
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

	product, err := models.GetProductByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "APPLE", product.Name)
	assert.Equal(t, 99.0, product.Price)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
}

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
		WithArgs("", -1.0). // Invalid data
		WillReturnError(errors.New("Validation error"))
	mock.ExpectRollback()

	product := &models.Product{
		Name:  "",
		Price: -1.0,
	}

	_, err = models.CreateProduct(product)
	assert.Error(t, err)
	assert.Equal(t, "Validation error", err.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
}

func TestDeleteProductWithInvalidData(t *testing.T) {
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
		WithArgs(999).
		WillReturnError(errors.New("Delete failed"))
	mock.ExpectRollback()

	rowsAffected, err := models.DeleteProduct(999)
	assert.Error(t, err)
	assert.Equal(t, "Delete failed", err.Error())
	assert.Equal(t, 0, rowsAffected)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
}

func TestUpdateProductWithInvalidData(t *testing.T) {
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
			AddRow(1, "APPLE", 50.0))

	expectedUpdate := "UPDATE `products` SET `name`=?,`price`=? WHERE `id` = ?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedUpdate)).
		WithArgs("APPLE_UPDATED", 200.0, 1).
		WillReturnError(errors.New("Update failed"))
	mock.ExpectRollback()

	updatedProduct := &models.Product{
		Name:  "APPLE_UPDATED",
		Price: 200.0,
	}

	err = models.UpdateProduct(1, updatedProduct)
	assert.Error(t, err)
	assert.Equal(t, "Update failed", err.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
}

func TestUpdateProductNotFound(t *testing.T) {
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

	err = models.UpdateProduct(1, &models.Product{Name: "UpdatedName"})

	assert.Error(t, err)
	assert.Equal(t, "Product not found", err.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
}

func TestGetAllProductsFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery(`^SELECT \* FROM ` + "`products`").
		WillReturnError(errors.New("Product not found"))

	models.DB = gormDB

	products, err := models.GetAllProducts()
	assert.Error(t, err)
	assert.Len(t, products, 0)
	assert.Equal(t, "Product not found", err.Error())
	assert.Nil(t, products)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
}

func TestGetProductByIDFailure(t *testing.T) {
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
		WillReturnError(errors.New("Product not found"))

	models.DB = gormDB

	product, err := models.GetProductByID(1)
	assert.Error(t, err)
	assert.Equal(t, "Product not found", err.Error())
	assert.Nil(t, product)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected: %s", err)
	}
}
