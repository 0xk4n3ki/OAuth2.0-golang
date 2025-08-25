package controllers

import (
	"github.com/OAuth2.0-golang/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
