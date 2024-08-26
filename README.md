# product-inventory-management-system
Develop a RESTful API service for a product inventory management system using Golang.    
     
## Overview
This is a simple RESTful API service for managing product inventory, 
developed in Golang using the MVC pattern. 
The system is built with the [Gin](https://github.com/gin-gonic/gin) framework 
and [GORM](https://gorm.io/) for ORM (Object-Relational Mapping). 
It allows users to add, update, delete, and retrieve product information. 
The product data is stored in a [MySQL](https://www.mysql.com/) database, 
ensuring performance and reliability with proper database handling techniques 
such as connection pooling, transactions, and error handling.    

---
## Features
* Create Product: Add a new product to the inventory.
* Retrieve Products: Get details of *all products* or a *specific product by ID.*
* Update Product: Modify details of an existing product.
* Delete Product: Remove a product from the inventory.

------
## API Endpoints
     
#### 1. Create Product
* Endpoint: POST /products
* Request Body:
```
{
  "name": "Product Name",
  "price": 100.0
}
```
* Response:
  * 201 Created: Product created successfully.
  * 400 Bad Request: Invalid input data.
  * 500 Internal Server Error: Database error.
     
#### 2. Retrieve All Products
* Endpoint: GET /products
* Response:
  * 200 OK: List of products.
  * 500 Internal Server Error: Database error.
      
#### 3. Retrieve Product by ID
* Endpoint: GET /products/{id}
* Response:
  * 200 OK: Product details.
  * 404 Not Found: Product not found.
  * 500 Internal Server Error: Database error.   
     
#### 4. Update Product
* Endpoint: PUT /products/{id}
* Request Body:
```
{
  "name": "Updated Product Name",
  "price": 120.0
}
```
#### 5. Delete Product
* Endpoint: DELETE /products/{id}
* Response:
  * 200 OK: Product deleted successfully.
  * 404 Not Found: Product not found.
  * 500 Internal Server Error: Database error.
