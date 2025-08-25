package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/0xk4n3ki/OAuth2.0-golang/config"
	"github.com/0xk4n3ki/OAuth2.0-golang/database"
	"github.com/0xk4n3ki/OAuth2.0-golang/helpers"
	"github.com/0xk4n3ki/OAuth2.0-golang/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var Validate = validator.New()

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("login")
		ctx.Redirect(http.StatusSeeOther, url)
	}
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("signup")
		ctx.Redirect(http.StatusSeeOther, url)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
		defer cancel()

	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GoogleCallback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state := ctx.Query("state")
		code := ctx.Query("code")

		token, err := config.AppConfig.GoogleLoginConfig.Exchange(ctx.Request.Context(), code)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Code-Token exchange failed"})
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
			return
		}
		defer resp.Body.Close()

		userData, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Json Parsing failed"})
			return
		}

		var user models.GoogleUser
		if err := json.Unmarshal(userData, &user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
			return
		}

		err = Validate.Struct(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "missing data returned by google"})
		}

		switch state {
		case "signup":
			helpers.AddUser(ctx, user)
		case "login":
			helpers.LoginUser(ctx, user)
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid State"})
		}
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken := ctx.GetHeader("refresh_token")
		if refreshToken == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
			return
		}

		claims, msg := helpers.ValidateToken(refreshToken)
		if msg != "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		newToken, newRefreshToken, _ := helpers.GenerateAllTokens(
			claims.Email,
			claims.First_name,
			claims.Last_name,
			claims.Uid,
		)
		if err := helpers.UpdateAllTokens(newToken, newRefreshToken, claims.Uid); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":         newToken,
			"refresh_token": newRefreshToken,
		})
	}
}
