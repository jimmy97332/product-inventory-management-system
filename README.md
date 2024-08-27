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

## Features
* Login: Authenticate users to generate a JWT token. Use this token to access other protected API endpoints.
     * The others are protected routes, meaning a valid JWT token must be provided in the request headers.
* Home: This endpoint welcomes authenticated users to the Product API.
* Create Product: Add a new product to the inventory.
* Retrieve Products: Get details of *all products* or a *specific product by ID.*
* Update Product: Modify details of an existing product.
* Delete Product: Remove a product from the inventory.

## API Endpoints

#### 1. Login
* Endpoint: POST /login/:user
* Response:
```
{
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzI0NzczNDEyfQ.rADyiQUIKj_nhVePIVVKOE0YHItcErz2Df_d9sL5sKI"
}
```
* Authorization:    
　* Mechanism: JWT (JSON Web Token)　　　
　* Middleware: The APIs under the /protected group are secured using JWT authentication.　You must include a valid JWT token in the request header as a Bearer token.　　　　　　
* Request Header Example:
```
"Authorization": "<your-jwt-token>"
```
#### 2. Home
* Endpoint: GET /protected/
* Response:
```
{
  "Welcome to the Product API"
}

```
#### 3. Create Product
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
     
#### 4. Retrieve All Products
* Endpoint: GET /products
* Response:
  * 200 OK: List of products.
  * 500 Internal Server Error: Database error.
      
#### 5. Retrieve Product by ID
* Endpoint: GET /products/{id}
* Response:
  * 200 OK: Product details.
  * 404 Not Found: Product not found.
  * 500 Internal Server Error: Database error.   
     
#### 6. Update Product
* Endpoint: PUT /products/{id}
* Request Body:
```
{
  "name": "Updated Product Name",
  "price": 120.0
}
```
#### 7. Delete Product
* Endpoint: DELETE /products/{id}
* Response:
  * 200 OK: Product deleted successfully.
  * 404 Not Found: Product not found.
  * 500 Internal Server Error: Database error.
 
## Database Schema
* Table Name: `products`
* Columns:
  * id: Integer (Primary Key, Auto Increment)
  * name: String
  * price: Float
* Sample SQL Script
```
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price FLOAT NOT NULL
);
```
* Model Migration
In the Go application, 
the Product model is automatically migrated to the database schema 
using GORM's AutoMigrate method. 
This ensures that the products table is created and kept up-to-date 
based on the current Product model definition.
```
db.AutoMigrate(&Product{})
```
This migration process handles creating the table 
if it doesn’t exist and updating the schema as necessary, 
based on changes to the Product model.
　　　　
##　Installation and Setup
###　Prerequisites
* Golang installed
* MySQL database    
###　Steps
####　Clone the Repository:
```
git clone https://github.com/jimmy97332/product-inventory.git
```
####　Install Dependencies:
* Before running the application, make sure to tidy up your Go module dependencies. This will add any missing dependencies and remove any that are no longer needed.　　　　
Run the following command:
```
go mod tidy
```
####　Configure Database:
* Update the database configuration in environment（.env） variables.
```
DATABASE_URL=root:password@tcp(localhost:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local
```
* Apply the SQL schema using the provided script.
####　Run the Application:
```
go run \main.go
```
    
#### Run Unit Tests
* To run unit tests for the API endpoints and database interactions, use the following command:
```
go test -v <file_name>
```
#### Coverage Report
* To see the overall coverage rate, use:
```
go test -v -cover
```
* To generate a test coverage report, use:
```
go test -coverprofile=coverage
go tool cover -html=coverage
```
The test coverage is approximately **80%**, ensuring a high level of code coverage for the project.
    
## Design Choices and Assumptions
* Web Framework: **Gin** framework is known for its speed and simplicity, which makes it an excellent choice for building RESTful APIs.
* Authentication: The project also provides middleware support, which was utilized to implement **JWT** authentication for protected routes.
* Database Handling: We use **GORM** as the ORM for interacting with the mySQL database, with connection pooling and safe transactions implemented.
* Error Handling: Proper error handling is implemented to ensure meaningful responses and uses the **Logrus** library which provides detailed logs that help in debugging and monitoring the application's behavior.
* JSON: All data between the client and server is exchanged in JSON format for simplicity and consistency.
Performance Considerations: Efficient database queries and connection pooling are used to handle performance concerns.
## Future Enhancements
* Expand the authentication and authorization system by integrating user accounts stored in the database to validate user credentials, in addition to the existing JWT-based authentication.
* Add pagination to the product listing endpoint.
* Include more comprehensive validation for input data.
