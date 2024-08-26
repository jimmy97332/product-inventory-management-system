package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Product struct {
	ID    int     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name  string  `json:"name" gorm:"column:name"`
	Price float64 `json:"price" gorm:"column:price"`
}

// init database
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Product{})

	// global DB
	DB = db
	return db, nil
}

func GetAllProducts() ([]Product, error) {
	var products []Product
	if err := DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func GetProductByID(id int) (*Product, error) {
	var product Product
	if err := DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func CreateProduct(product *Product) (int, error) {
	if err := DB.Create(product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

func DeleteProduct(id int) (int, error) {
	if DB == nil {
		return 0, DB.Error
	}

	var product Product
	result := DB.Delete(&product, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func UpdateProduct(id uint, updatedData *Product) error {
	var product Product
	if err := DB.First(&product, id).Error; err != nil {
		return err
	}

	if updatedData.Name != "" {
		product.Name = updatedData.Name
	}
	if updatedData.Price != 0 {
		product.Price = updatedData.Price
	}

	if err := DB.Save(&product).Error; err != nil {
		return err
	}
	return nil
}
