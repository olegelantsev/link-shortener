package api

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

func VerifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized!"))
					if err != nil {
						return nil, err
					}
				}
				return "", nil

			})
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}
			}
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	})
}
