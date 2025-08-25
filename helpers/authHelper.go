package helpers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/0xk4n3ki/OAuth2.0-golang/database"
	"github.com/0xk4n3ki/OAuth2.0-golang/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(ctx *gin.Context, gUser models.GoogleUser) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	count, err := database.UserCollection.CountDocuments(c, bson.M{"email": gUser.Email})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking user"})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	var user models.User
	user.ID = primitive.NewObjectID()
	user.First_name = &gUser.Given_name
	user.Last_name = &gUser.Family_name
	user.Created_at = time.Now().UTC()
	user.Updated_at = time.Now().UTC()
	user.User_id = user.ID.Hex()

	token, refreshToken, _ := GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	_, insertErr := database.UserCollection.InsertOne(c, user)
	if insertErr != nil {
		log.Println("error inserting user: ", insertErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not created"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user_id": user.User_id})
}

func LoginUser(ctx *gin.Context, user models.GoogleUser) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	var foundUser models.User
	err := database.UserCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User don't exist"})
		return
	}

	token, refreshToken, _ := GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *&foundUser.User_id)

	if err := UpdateAllTokens(token, refreshToken, foundUser.User_id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tokens"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":       foundUser.User_id,
		"email":         foundUser.Email,
		"first_name":    foundUser.First_name,
		"last_name":     foundUser.Last_name,
		"token":         foundUser.Token,
		"refresh_token": foundUser.Refresh_token,
	})
}
