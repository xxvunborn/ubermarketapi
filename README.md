Documentation for UberMarket API
===

Api for ubermarket writen in go

Routes
---

 For users

```
[GET]    	/api/v1/users
[GET]    	/api/v1/user/:id
[POST]  	/api/v1/users
[PUT]    	/api/v1/user/:id
[DELETE] 	/api/v1/user/:id
```
 For Products

```
[GET]		/api/products
[GET] 		/api/product/:id
```

For Authenticated

```
[POST] /api/v1/authenticated
```

```
Mode of use:
send POST with json params

Example:
{
	"email": "unborn.system@gmail.com",
	"password":"12333456"
}

return
 	"token": "Authenticated" 
return
	 "error" : "User and/or password is incorrect" 
```

Dependencies
--
`github.com/gin-gonic/gin` For manage roues

`github.com/go-sql-driver/mysql` For manage sql handlers

`gopkg.in/gorp.v1` For manage persistence in DB

How to install
--
 Install Depencencies


```
go get github.com/gin-gonic/gin
go get github.com/go-sql-driver/mysql
go gopkg.in/gorp.v1
```
How to run API
--
 run in comand line

```
cd ~/API
go run main.go
``` 

 Go to browser

```
localhost:3000
```

TO DO
--
* Fix NULL String error for NULL values in tables

* Add other routes 
