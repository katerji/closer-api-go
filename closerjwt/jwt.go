package closerjwt

import (
	"closer-api-go/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

func VerifyToken(token string) (model.User, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		return model.User{}, errors.New("error parsing token")
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		jsonClaims, err := json.Marshal(claims)
		if err != nil {
			return model.User{}, errors.New("error parsing token")
		}
		var user model.User
		if err := json.Unmarshal(jsonClaims, &user); err != nil {
			return model.User{}, errors.New("error parsing token")
		}
		return user, nil
	}
	return model.User{}, errors.New("invalid token")

}

func CreateJwt(user model.User) (string, error) {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	fmt.Println()
	//	return "", err
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":           user.Id,
		"phone_number": user.PhoneNumber,
		"name":         user.Name,
	})
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}
