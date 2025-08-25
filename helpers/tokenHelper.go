package helpers

import (
	"context"
	"os"
	"time"

	"github.com/0xk4n3ki/OAuth2.0-golang/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, uid string) (token, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Minute).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func UpdateAllTokens(token, refreshToken, uid string) error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateObj := bson.M{
		"token":         token,
		"refresh_token": refreshToken,
		"updated_at":    time.Now().UTC(),
	}

	filter := bson.M{"user_id": uid}
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := database.UserCollection.UpdateOne(
		c,
		filter,
		bson.M{"$set": updateObj},
		&opts,
	)

	return err
}

func ValidateToken(signedToken string) (*SignedDetails, string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err.Error()
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		return nil, "Token is Invalid"
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, "Token is Expired"
	}

	return claims, ""
}
