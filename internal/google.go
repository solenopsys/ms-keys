package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type GoogleAuth struct {
	Certs map[string]string
}

type CertsResponse struct {
	Keys []struct {
		KeyId string   `json:"keyId"`
		X5c   []string `json:"x5c"`
	} `json:"keys"`
}

func (g *GoogleAuth) Init() {

	const certsURL = "https://www.googleapis.com/oauth2/v1/certs"

	g.Certs = make(map[string]string)

	resp, err := http.Get(certsURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var certs CertsResponse
	err = json.Unmarshal(body, &certs)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cert := range certs.Keys {
		g.Certs[cert.KeyId] = cert.X5c[0]
	}

}

func checkToken(tokenString string) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the algorithm used for signing the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Use the secret key for verifying the token
		return []byte("secret"), nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		// Token is valid, you can access the claims
		fmt.Printf("Subject: %s\n", claims.Subject)
		fmt.Printf("Name: %s\n", claims.VerifyAudience("", true))
		fmt.Printf("Issued At: %v\n", claims.IssuedAt)
		fmt.Printf("Expiry: %v\n", claims.ExpiresAt)
	} else {
		fmt.Println("Invalid token")
	}

}
