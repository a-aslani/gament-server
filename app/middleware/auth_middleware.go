package middleware

import (
	"Server/app/constants"
	"Server/app/response"
	"Server/app/utility"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	HeaderAuthorization = "Authorization"
	HeaderRegister      = "Register"
)

//Check RS256 token
func UserAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader(HeaderAuthorization))
		if header == "" {
			c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.UnauthorizedError{Error: "ارسال توکن شناسایی الزامی میباشد"}})
			c.Abort()
		} else {
			token, err := jwt.Parse(c.GetHeader(HeaderAuthorization), func(token *jwt.Token) (interface{}, error) {
				return utility.VerifyKey, nil
			})
			if err == nil {
				if token.Valid {
					//Send claims for receive in routes
					c.Set(constants.TokenClaims, token.Claims)
					c.Next()
				}
			} else {
				c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.UnauthorizedError{Error: "توکن شناسایی نامعتبر است"}})
				c.Abort()
			}
		}
	}
}

//Check RS256 token for admin or owner
func AdminAndOwnerAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader(HeaderAuthorization))
		if header == "" {
			c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.UnauthorizedError{Error: "ارسال توکن شناسایی الزامی میباشد"}})
			c.Abort()
		} else {
			token, err := jwt.Parse(c.GetHeader(HeaderAuthorization), func(token *jwt.Token) (interface{}, error) {
				return utility.VerifyKey, nil
			})
			if err == nil {
				if token.Valid {
					claims := token.Claims.(jwt.MapClaims)
					role := claims["role"].(string)
					if role == constants.AdminRole {
						c.Set(constants.TokenClaims, token.Claims)
						c.Next()
					} else if role == constants.OwnerRole {
						c.Set(constants.TokenClaims, token.Claims)
						c.Next()
					} else {
						c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.UnauthorizedError{Error: "شما اجازه ی دسترسی به این بخش را ندارید"}})
						c.Abort()
					}
				}
			} else {
				c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.UnauthorizedError{Error: "توکن شناسایی نامعتبر است"}})
				c.Abort()
			}
		}
	}
}

//Check HS256 token
func RegisterAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := strings.TrimSpace(c.GetHeader(HeaderRegister))

		if header == "" {
			c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.UnauthorizedError{Error: "ارسال توکن شناسایی الزامی میباشد"}})
			c.Abort()
		} else {
			token, err := jwt.Parse(c.GetHeader(HeaderRegister), func(token *jwt.Token) (interface{}, error) {
				return []byte(constants.SecretKey), nil
			})
			if err == nil {
				if token.Valid {
					//Send claims for receive in routes
					c.Set(constants.TokenClaims, token.Claims)
					c.Next()
				}
			} else {
				c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.UnauthorizedError{Error: "توکن شناسایی نامعتبر است"}})
				c.Abort()
			}
		}

	}
}
