package hjwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"strings"
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

func Validate(secret []byte, tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
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

func Generate(secret []byte, uid string, defaultRole string, roles []string, user interface{}, extra map[string]interface{}, exp time.Duration) (string, error) {
	expirationTime := time.Now().Add(exp)
	hasura := make(map[string]interface{})
	hasura[`x-hasura-allowed-roles`] = []string{defaultRole}
	hasura[`x-hasura-default-role`] = defaultRole
	hasura[`x-hasura-user-id`] = uid
	roles = append(roles, defaultRole)
	hasura[`x-hasura-user-roles`] = fmt.Sprintf(`{%s}`, strings.Join(roles, `,`))
	for k, v := range extra {
		vt := reflect.TypeOf(v).String()
		switch vt {
		case "[]string":
			hasura[fmt.Sprintf("x-hasura-user-%s", k)] = fmt.Sprintf(`{%s}`, strings.Join(v.([]string), `,`))
		default:
			hasura[fmt.Sprintf("x-hasura-user-%s", k)] = v
		}
	}

	claims := &Claims{
		Uid:    uid,
		Role:   defaultRole,
		User:   user,
		Hasura: hasura,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "HASURA-JWT-GENERATOR",
			IssuedAt:  time.Now().Add(-1 * time.Minute).Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
