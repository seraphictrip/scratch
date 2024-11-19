# DESIGN



## Use cases

A user should be able to register for the site

A user should be able to login the site

A user should be able to get a list of products

A user should be able to get info about self

## Non-functional
Databases can be stood up, torn down simply for each developer env.

## Infrastructure

Golang for API Server

SQL for database

Docker to containerize 

## High level design

POST /register

POST /login

GET /products

GET /users/{userId}


## Entities

Users
- ID
- firstName
- lastName
- email
- password

Products
- ID
- name
- description
- image
- quantity
- createdAt

Orders
- ID
- userId
- total
- status
- address
- createdAt

Orders_Items
- ID
- orderId
- productID
- quantity
- price