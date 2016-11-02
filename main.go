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
	//dbmap.AddTableWithName(User{}, "User").SetKeys(true, "Id")
	//err = dbmap.CreateTablesIfNotExists()
	//checkErr(err, "Create table failed")

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
	err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)

	if err == nil {
		var json User
		c.Bind(&json)

		user_id, _ := strconv.ParseInt(id, 0, 64)

		user := User{
			Id:       user_id,
			Email:    json.Email,
			Password: json.Password,
		}

		if user.Email != "" && user.Password != "" {
			_, err = dbmap.Update(&user)

			if err == nil {
				c.JSON(200, user)
			} else {
				checkErr(err, "Update failed")
			}
		} else {
			c.JSON(422, gin.H{"error": "Field are empty"})
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
