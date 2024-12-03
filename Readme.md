# IMS-Zedeks API Documentation

## Products Endpoints

### Get All Products
- **URL**: `/products`
- **Method**: `GET`
- **Success Response**:
  - **Code**: 200
  - **Content**: Array of product objects
- **No Content Response**:
  - **Code**: 204
  - **Content**: Empty array

### Create Product
- **URL**: `/products`
- **Method**: `POST`
- **Data Params**:
  ```json
  {
    "name": "string",
    "category_id": "uuid",
    "price": "float64",
    "quantity": "integer",
    "image_url": "string (optional)",
    "supplier_id": "uuid"
  }
  ```
- **Success Response**:
  - **Code**: 201
  - **Content**: Created product object
- **Error Responses**:
  - **Code**: 400
    - **Content**: `{"error": "Invalid request body"}`
  - **Code**: 500
    - **Content**: `{"error": "Failed to create product"}`

### Get Single Product
- **URL**: `/products/:id`
- **Method**: `GET`
- **URL Params**: `id=[uuid]`
- **Success Response**:
  - **Code**: 200
  - **Content**: Product object
- **Error Response**:
  - **Code**: 404
    - **Content**: `{"error": "Product not found"}`

### Update Product
- **URL**: `/products/:id`
- **Method**: `PUT`
- **URL Params**: `id=[uuid]`
- **Data Params**: Same as Create Product
- **Success Response**:
  - **Code**: 200
  - **Content**: Updated product object
- **Error Responses**:
  - **Code**: 404
    - **Content**: `{"error": "Product not found"}`
  - **Code**: 400
    - **Content**: `{"error": "Invalid request body"}`
  - **Code**: 500
    - **Content**: `{"error": "Failed to update product"}`

### Delete Product
- **URL**: `/products/:id`
- **Method**: `DELETE`
- **URL Params**: `id=[uuid]`
- **Success Response**:
  - **Code**: 204
  - **Content**: No Content
- **Error Responses**:
  - **Code**: 404
    - **Content**: `{"error": "Product not found"}`
  - **Code**: 500
    - **Content**: `{"error": "Failed to delete product"}`

### Get Products by Category
- **URL**: `/categories/:categoryId/products`
- **Method**: `GET`
- **URL Params**: `categoryId=[uuid]`
- **Success Response**:
  - **Code**: 200
  - **Content**: Array of product objects
- **Error Responses**:
  - **Code**: 400
    - **Content**: `{"error": "Invalid category ID format"}`
  - **Code**: 500
    - **Content**: `{"error": "Failed to fetch products"}`

### Get Products by Supplier
- **URL**: `/suppliers/:supplierId/products`
- **Method**: `GET`
- **URL Params**: `supplierId=[uuid]`
- **Success Response**:
  - **Code**: 200
  - **Content**: Array of product objects
- **Error Responses**:
  - **Code**: 400
    - **Content**: `{"error": "Invalid supplier ID format"}`
  - **Code**: 500
    - **Content**: `{"error": "Failed to fetch products"}`

## Categories Endpoints

### Get All Categories
- **URL**: `/categories`
- **Method**: `GET`

### Create Category
- **URL**: `/categories`
- **Method**: `POST`

### Get Single Category
- **URL**: `/categories/:id`
- **Method**: `GET`

### Update Category
- **URL**: `/categories/:id`
- **Method**: `PUT`

### Delete Category
- **URL**: `/categories/:id`
- **Method**: `DELETE`

## Suppliers Endpoints

### Get All Suppliers
- **URL**: `/suppliers`
- **Method**: `GET`

### Create Supplier
- **URL**: `/suppliers`
- **Method**: `POST`

### Get Single Supplier
- **URL**: `/suppliers/:id`
- **Method**: `GET`

### Update Supplier
- **URL**: `/suppliers/:id`
- **Method**: `PUT`

### Delete Supplier
- **URL**: `/suppliers/:id`
- **Method**: `DELETE`

