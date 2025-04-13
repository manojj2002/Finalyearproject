package controllers

import (
	"CVC_ragh/config" // Adjust the path to your actual module name
	"CVC_ragh/models" // For the User struct
	"CVC_ragh/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	clientID     = "Ov23liIqM2f5hgTtVp32"
	clientSecret = "e1ae9c413c43309dfe372ee640906d1e6064f366"
	oauthURL     = "https://github.com/login/oauth"
	apiURL       = "https://api.github.com/user"
	redirectURI  = "http://localhost:4000/api/auth/github/callback"
)

func HandleGitHubLogin(c *gin.Context) {
	url := fmt.Sprintf("%s/authorize?client_id=%s&redirect_uri=%s&scope=user", oauthURL, clientID, redirectURI)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGitHubCallback(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Code Provided"})
		return
	}

	token, err := getAccessToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	userData, err := getUserData(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get User data"})
		return
	}

	username := userData["login"].(string)
	email := ""
	if val, ok := userData["email"].(string); ok && val != "" {
		email = val
		fmt.Println(email)
	} else {
		fmt.Println("⚠️  Email not available in GitHub profile.")

	}
	//fmt.Println(userData)
	// Check or create user in MongoDB
	userId, err := findOrCreateUserInMongo(username, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}

	// Issue your own JWT
	myJWT, err := utils.CreateJWT(username, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT"})
		return
	}
	fmt.Println(myJWT)
	redirectURL := "http://localhost:5173/githubLogin?token=" + myJWT
	c.Redirect(http.StatusSeeOther, redirectURL)

}

func getAccessToken(code string) (string, error) {
	url := fmt.Sprintf("%s/access_token", oauthURL)
	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s", clientID, clientSecret, code, redirectURI)

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	values, err := urlValuesFromBody(string(body))
	if err != nil {
		return "", err
	}

	return values["access_token"], nil
}

func getUserData(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, err
}

func urlValuesFromBody(body string) (map[string]string, error) {
	values := make(map[string]string)
	for _, pair := range strings.Split(body, "&") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid body format")
		}
		values[kv[0]] = kv[1]
	}
	return values, nil
}
func findOrCreateUserInMongo(username string, email string) (string, error) {
	ctx := context.TODO()
	collection := config.GetDB().Collection("users")

	var user models.User

	// Try to find the user
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found, so create a new one
			newUser := models.User{
				Username: username,
				Email:    email,
			}

			res, insertErr := collection.InsertOne(ctx, newUser)
			if insertErr != nil {
				return "", fmt.Errorf("failed to insert new user: %v", insertErr)
			}

			oid, ok := res.InsertedID.(primitive.ObjectID)
			if !ok {
				return "", fmt.Errorf("failed to parse inserted ID as ObjectID")
			}
			return oid.Hex(), nil
		}
		// Some other DB error
		return "", fmt.Errorf("failed to query user: %v", err)
	}

	// User exists, return their ID
	return user.ID.Hex(), nil
}
