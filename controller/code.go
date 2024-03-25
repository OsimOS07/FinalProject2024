package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"test/app/structs"
	"time"

	// "github.com/bytedance/sonic/option"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SignUp(c *gin.Context) {
	var shablon structs.Signup
	c.ShouldBindJSON(&shablon)

	if shablon.Email == "" || shablon.Login == "" || shablon.Password == "" {
		c.JSON(400, "Empty fill")
	} else {

		CTXReturned, ClientReturned := Connection()

		connection := ClientReturned.Database("OsimDB").Collection("Users")

		SingleResult := connection.FindOne(CTXReturned, bson.M{
			"email": shablon.Email,
		})
		var shablons structs.Signup
		SingleResult.Decode(&shablons)

		if shablons.Email != "" {
			c.JSON(404, "Erorr emil exist")
		} else {
			GeneratedId := primitive.NewObjectID().Hex()

			HashedPassword := HashPass(shablon.Password)

			connection.InsertOne(CTXReturned, bson.M{
				"_id":      GeneratedId,
				"login":    shablon.Login,
				"password": HashedPassword,
				"email":    shablon.Email,
			})
		}
	}
}
func Connection() (context.Context, *mongo.Client) {
	url := options.Client().ApplyURI("mongodb://192.168.43.246:27018")
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	Client, err := mongo.Connect(ctx, url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return ctx, Client

}
func HashPass(password string) string {
	HashedPssword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(HashedPssword)
}

func ComPass(DBPAssword string, UserPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(DBPAssword), []byte(UserPassword))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return err == nil
}

func Login(c *gin.Context) {
	var sablon structs.Login
	c.ShouldBindJSON(&sablon)

	if sablon.Login == "" || sablon.Password == "" {
		c.JSON(404, "Erorr")
	} else {
		Ctx, CLient := Connection()
		connect := CLient.Database("OsimDB").Collection("Users")
		result := connect.FindOne(Ctx, bson.M{
			"login": sablon.Login,
		})
		var userData structs.Login
		result.Decode(&userData)
		fmt.Printf("userData: %v\n", userData)
		IsValid := ComPass(userData.Password, sablon.Password)
		if IsValid {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:  "Qamch",
				Value: userData.Id,
			})
			c.JSON(200, "Sucess")
		} else {

			c.JSON(401, "Erorr Invalid Password or login")
		}
		fmt.Printf("userData: %v\n", userData)

	}
}
func AddProduct(c *gin.Context) {
	cookieData, err := c.Request.Cookie("Qamch")

	fmt.Printf("cookieData: %v\n", cookieData)
	if err != nil {

		c.JSON(400, "erorr")
	} else {
		var shablon structs.AddProduct
		c.ShouldBindJSON(&shablon)

		if shablon.ProductName == "" || shablon.ProductPrice == 0 || shablon.ShortDescription == "" {
			c.JSON(404, "Emtpy field")
		} else {

			bytedata, err := base64.StdEncoding.DecodeString(shablon.Image)
			if err != nil {
				fmt.Printf("err: %v\n", err)

			}
			rand := rand.Intn(10000)
			name := fmt.Sprintf("./Project/%v.jpg", rand)
			erorr := ioutil.WriteFile(name, bytedata, 0644)

			if erorr != nil {
				fmt.Printf("erorr: %v\n", erorr)
			}
			Ctx, product := Connection()
			connect := product.Database("OsimDB").Collection("Products")
			img := fmt.Sprintf("%v.png", rand)
			connect.InsertOne(Ctx, bson.M{
				"_id":              primitive.NewObjectID().Hex(),
				"productname":      shablon.ProductName,
				"productprice":     shablon.ProductPrice,
				"shortdescription": shablon.ShortDescription,
				"image":            img,
				"ownerid":          cookieData.Value,
			})
		}
		fmt.Printf("shablon: %v\n", shablon)
	}
}

func List(c *gin.Context) {
	var slice = []structs.AddProduct{}
	cookieData, err := c.Request.Cookie("Qamch")

	fmt.Printf("cookieData: %v\n", cookieData)
	if err != nil {

		c.JSON(400, "erorr")
	} else {

		CTX, Getlist := Connection()
		connect := Getlist.Database("OsimDB").Collection("Products")
		Result, _ := connect.Find(CTX, bson.M{
			"ownerid": cookieData.Value,
		})
		for Result.Next(CTX) {
			var shablon structs.AddProduct
			Result.Decode(&shablon)
			slice = append(slice, shablon)
		}
		c.JSON(200, slice)

	}

}
func Delete(c *gin.Context) {
	cookieData, err := c.Request.Cookie("Qamch")

	fmt.Printf("cookieData: %v\n", cookieData)
	if err != nil {

		c.JSON(400, "erorr")
	} else {
		var DeletTempt structs.Delete
		c.ShouldBindJSON(&DeletTempt)
		CTX, Connect := Connection()

		connect := Connect.Database("OsimDB").Collection("Products")
		connect.DeleteOne(CTX, bson.M{
			"_id": DeletTempt.ItemID,
		})
		c.JSON(200, "Deleted")

	}
}
