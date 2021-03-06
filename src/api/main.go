package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"strconv"
)

type User struct {
	Id             int64  `db:"id" json:"id"`
	Email          string `db:"email" json:"email"`
	Password       string `db:"password" json:"password"`
	Name           string `db:"name" json:"name"`
	Phone          string `db:"phone" json:"phone"`
	ConfirmedEmail string `db:"confirmed_email" json:"confirmed_email"`
}

type Product struct {
	Id             int64  `db:id json:"id"`
	Name           string `db:"name" json:"name"`
	Brand          string `db:"brand" json:"brand"`
	Description    string `db:"description" json:"description"`
	Price          int    `db:"price" json:"price"`
	Stock          int    `db:"stock" json:"stock"`
	AvailableStock int    `db:"available_stock" json:"available_stock"`
	IdCategory     int    `db:"id_category" json:"id_category"`
}

type Order struct {
	Id          int64 `db:id json"id"`
	IdUser      int64 `db:"id_user" json:"id_user"`
	OrderNumber int64 `db:"order_number" json:"order_number"`
	IdProduct   int64 `db:"id_product" json:"id_product"`
}

type Authenticator struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type Client struct {
	Id           int64  `db:id json:"id"`
	IdUser       int64  `db:"id_user" json:"id_user"`
	DueDate      int64  `db:"due_date" json:"due_date"`
	NumberOfCard int64  `db:"number_of_card" json:"number_of_card"`
	Address      string `db:"address" json:"address"`
	TypeOfCard   string `db:"type_of_card" json:"type_of_card"`
}

var dbmap = initDb()

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:1234@/superuber")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	// Create table if not exists
	dbmap.AddTableWithName(User{}, "User").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")

	return dbmap
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		v1.GET("/users", getUsers)
		v1.GET("/user/:id", getUser)
		v1.POST("/users", postUser)
		v1.PUT("/users/:id", updateUser)
		v1.DELETE("/users/:id", deleteUser)

		v1.GET("/products", getProducts)
		v1.GET("/product/:id", getProduct)

		v1.POST("/authentication", authentication)

		v1.POST("/order", postOrder)
		v1.POST("/order/verify/:order", verifyOrder)
	}
	r.Run(":3000")
}

func getUsers(c *gin.Context) {
	var users []User
	_, err := dbmap.Select(&users, "SELECT * FROM user")

	if err == nil {
		c.JSON(200, users)
	} else {
		c.JSON(404, err)
	}
	// curl -i http://localhost:3000/api/v1/users
}

func getUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

	if err == nil {
		user_id, _ := strconv.ParseInt(id, 0, 64)

		content := &User{
			Id:       user_id,
			Email:    user.Email,
			Password: user.Password,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
	// curl -i http://localhost:3000/api/v1/users/1
}

func postUser(c *gin.Context) {
	var user User
	c.Bind(&user)

	if user.Email != "" && user.Password != "" {

		if insert, _ := dbmap.Exec(`INSERT INTO user (email, password) VALUES (?, ?)`, user.Email, user.Password); insert != nil {
			user_id, err := insert.LastInsertId()
			if err == nil {
				content := &User{
					Id:       user_id,
					Email:    user.Email,
					Password: user.Password,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}
	} else {
		c.JSON(422, gin.H{"error": "field are empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func updateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

	if err == nil {
		var json User
		c.Bind(&json)

		user_id, _ := strconv.ParseInt(id, 0, 64)

		user := User{
			Id:       user_id,
			Email:    json.Email,
			Password: json.Password,
			Name:     json.Name,
		}

		if user.Email != "" && user.Password != "" {
			_, err = dbmap.Update(&user)

			if err == nil {
				c.JSON(200, user)
			} else {
				checkErr(err, "Update failed")
			}
		} else {
			c.JSON(422, gin.H{"error": "No data"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func deleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT id FROM User WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&user)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: " deleted"})
		} else {
			checkErr(err, "Deleted field")
		}
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}

func getProducts(c *gin.Context) {
	var products []Product
	_, err := dbmap.Select(&products, "SELECT * FROM product")

	if err == nil {
		c.JSON(200, products)
	} else {
		c.JSON(404, err)
	}
	// curl -i http://localhost:3000/api/v1/users
}

func getProduct(c *gin.Context) {
	id := c.Params.ByName("id")
	var product Product
	err := dbmap.SelectOne(&product, "SELECT * FROM product WHERE id=?", id)

	if err == nil {
		product_id, _ := strconv.ParseInt(id, 0, 64)

		content := &Product{
			Id:             product_id,
			Name:           product.Name,
			Brand:          product.Brand,
			Description:    product.Description,
			Price:          product.Price,
			Stock:          product.Stock,
			AvailableStock: product.AvailableStock,
			IdCategory:     product.IdCategory,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
	// curl -i http://localhost:3000/api/v1/users/1
}
func authentication(c *gin.Context) {
	var user User
	c.Bind(&user)

	if user.Email != "" && user.Password != "" {
		err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE email=? and password=?", user.Email, user.Password)
		if err == nil {
			c.JSON(200, gin.H{"token": "authenticated"})
		} else {
			c.JSON(404, gin.H{"error": "User and/or password is incorrect"})
		}

	}
}

func postOrder(c *gin.Context) {
	var order Order
	c.Bind(&order)

	if order.IdProduct != 0 {

		if insert, _ := dbmap.Exec(`INSERT INTO order_Product (id_user, id_product, order_number) VALUES (?, ?, ?)`, order.IdUser, order.IdProduct, order.OrderNumber); insert != nil {
			order_id, err := insert.LastInsertId()
			if err == nil {
				content := &Order{
					Id:          order_id,
					IdUser:      order.IdUser,
					IdProduct:   order.IdProduct,
					OrderNumber: order.OrderNumber,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}
	} else {
		c.JSON(422, gin.H{"error": "field are empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func verifyOrder(c *gin.Context) {
	order := c.Params.ByName("order")
	err := dbmap.SelectOne("SELECT * FROM order_product WHERE id=?", order)

	if err != nil {
		c.JSON(200, gin.H{"verify": "Verify Complete"})
		// Change delivered 0 to 1 in purchase table  to end the circle
		// dbmap.Exec(`UPDATE purchase SET delivered=0 where id_order = ?`, order)
	}
	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}
