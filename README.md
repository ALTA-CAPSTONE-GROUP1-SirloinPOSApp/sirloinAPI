
### Build App & Database

![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)
![Visual Studio Code](https://img.shields.io/badge/Visual%20Studio%20Code-0078d7.svg?style=for-the-badge&logo=visual-studio-code&logoColor=white)
![MySQL](https://img.shields.io/badge/mysql-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![AWS](https://img.shields.io/badge/Amazon_AWS-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![Cloudflare](https://img.shields.io/badge/Cloudflare-F38020?style=for-the-badge&logo=Cloudflare&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Midtrans](https://img.shields.io/badge/Midtrans-FF6C37?style=for-the-badge&logo=midtrans&logoColor=white)
![Gmail](https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white)

# SIRLOIN POS API

This is a GO language REST API group project organized by Group 1. This API is used to run SILOIN POS applications. SIRLOIN POS is Point of Sale application with targeted user is small business like 'warung' etc.

This application has features as listed below. 


# Features
## User:
- Register
- Login
- Show profile
- Edit profile
- Deactive account
- Register Device

<div>

<details>

| Feature User | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| POST | /register | - | NO | Register new users (tenant). |
| POST | /login  | - | NO | Log in into tenant account.  |
| GET | /users | - | YES | Get tenant information details. |
| PUT | /users | - | YES | Edit tenant details. |
| DELETE | /users | - | YES | Delete/deactive account. |
| POST | /register_device | - | YES | Register device token for notification. |

</details>

<div>

## Admin :
- Get admin products
- Get admin selling history
- Get admin selling details

<div>

<details>

| Feature Product | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| GET | /products/admin | - | YES | Get admin products. |
| GET | /transactions/admin  | - | YES | Get admin selling history.  |
| GET | /transactions/{transaction_id}/admin | TRANSACTION ID | YES | Get admin selling details. |

</details>

</div>

## Product :
- Add product
- Show all product
- Edit product
- Show detail product
- Delete product

<div>

<details>

| Feature Product | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| POST | /products | - | YES | Add new product for user and admin. |
| GET | /products  | - | YES | Get all tenant products.  |
| PUT | /products | PRODUCT ID | YES | Edit tenant and admin product. |
| GET | /products | PRODUCT ID | YES | Get product details for tenant and admin. |
| DELETE | /products | PRODUCT ID | YES | Delete product for tenant and admin. |

</details>

</div>

## Customers :
- Add new customer
- Get all customer
- Edit customer details
- Get customer detail by ID

<div>

<details>

| Feature Cart | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| POST | /customers | - | YES | Register new customer. |
| GET | /customers  | - | YES | Get all tenant customers.  |
| PUT | /customers | CUSTOMER ID | YES | Edit customer detail. |
| GET | /customers  | CUSTOMER ID | YES | Get customer detail by ID.  |

</details>

</div>

## Transactions :
- Create new transaction sell
- Get selling or buying history
- Edit transaction status
- Get transaction detail
- Create new transaction buy

<div>

<details>

| Feature Cart | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| POST | /transactions | - | YES | Create new selling transaction. |
| GET | /transactions  | - | YES | Get buying or selling history depends on query param.  |
| PUT | /transactions | TRANSACTION ID | YES | Edit transaction status. |
| GET | /transactions | TRANSACTION ID | YES |  GET transaction details. |
| POST | /transactions | TRANSACTION ID | YES | Create new buying transaction
</details>

</div>


# ERD
![ERD](https://mediasosial.s3.ap-southeast-1.amazonaws.com/Sirloin.drawio.png "ERD")

# API Documentations

[Click here](https://app.swaggerhub.com/apis-docs/CAPSTONE-Group1/sirloinPOSAPI/1.0.0) to see documentations.


## How to Install To Your Local

- Clone it

```
$ git clone https://github.com/ALTA-CAPSTONE-GROUP1-SirloinPOSApp/sirloinAPI.git
```

- Go to directory

```
$ cd SirloinAPI
```
## Authors 👑
-  [![GitHub](https://img.shields.io/badge/ari-muhammad-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)](https://github.com/amrizal94)

-  [![GitHub](https://img.shields.io/badge/fauzan-putra-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)](https://github.com/mfauzanptra)


 <p align="right">(<a href="#top">back to top</a>)</p>
<h3>
<p align="center">:copyright: February 2023 </p>
</h3>
<!-- end -->
<!-- comment -->
