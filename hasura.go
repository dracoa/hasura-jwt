package hjwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Claims struct {
	Uid    string                 `json:"uid"`
	Role   string                 `json:"role"`
	User   interface{}            `json:"user"`
	Hasura map[string]interface{} `json:"https://hasura.io/jwt/claims"`
	jwt.StandardClaims
}

type User struct {
	Uid   string
	Roles []string
	Extra map[string]interface{}
}

var hasuraJwtSecret []byte

func init() {
	_ = godotenv.Load()
	hasuraJwtSecret = []byte(os.Getenv("HASURA_JWT_SIGN_KEY"))
}

func Validate(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hasuraJwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jsonString, _ := json.Marshal(claims)
		s := &Claims{}
		err = json.Unmarshal(jsonString, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	} else {
		return nil, errors.New("token claims not okay")
	}
}

func Generate(uid string, roles []string, user interface{}, extra map[string]interface{}, exp time.Duration) (string, error) {
	expirationTime := time.Now().Add(exp)
	if len(roles) == 0 {
		return "", errors.New("user should at least has one role")
	}
	hasura := make(map[string]interface{})
	hasura[`x-hasura-allowed-roles`] = roles
	hasura[`x-hasura-default-role`] = roles[0]
	hasura[`x-hasura-user-id`] = uid
	hasura[`x-hasura-user-roles`] = roles[0]
	for k, v := range extra {
		hasura[fmt.Sprintf("x-hasura-user-%s", k)] = v
	}

	claims := &Claims{
		Uid:    uid,
		Role:   roles[0],
		User:   user,
		Hasura: hasura,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "HASURA-JWT-GENERATOR",
			IssuedAt:  time.Now().Add(-1 * time.Minute).Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hasuraJwtSecret)
}
