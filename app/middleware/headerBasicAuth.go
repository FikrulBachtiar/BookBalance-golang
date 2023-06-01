package middleware

import (
	"bookbalance/app/configs"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetUser(db *sql.DB, issuer_code string, issuer_secret string) int {
	sqlQuery := fmt.Sprintf("SELECT third_party_code FROM mtr.t_mtr_third_party WHERE third_party_code = '%s' AND third_party_secret = '%s' AND active_status = 1 LIMIT 1", issuer_code, issuer_secret);
	var user string;

	err := db.QueryRow(sqlQuery).Scan(&user);
	if err != nil {
		return 1;
	}
	
	return 0;
}

func AuthBasic(db *sql.DB) echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			header := c.Request().Header.Get("Authorization");
			isBasic := strings.HasPrefix(header, "Basic ");
			
			if !isBasic {
				response := &configs.Response{
					Status: http.StatusUnauthorized,
					Code: 401,
					Message: "Unauthorized",
				}
				return response.ResponseMiddleware(c);
			}

			decode, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, "Basic "));
			if err != nil {
				response := &configs.Response{
					Status: http.StatusUnauthorized,
					Code: 401,
					Message: "Unauthorized",
				}
				return response.ResponseMiddleware(c);
			}

			authSplit := strings.Split(string(decode), ":");
			issuer_code := authSplit[0];
			issuer_secret := authSplit[1];

			code := GetUser(db, issuer_code, issuer_secret);
			if code != 0 {
				response := &configs.Response{
					Status: http.StatusUnauthorized,
					Code: 401,
					Message: "Unauthorized",
				}
				return response.ResponseMiddleware(c);
			}

			return next(c);
		}
	}
}