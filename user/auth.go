package user

import (
	config "PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Location string `form:"location" json:"location" binding:"required"`
}

var identityKey = "id"

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).userid,
		"text":     "Hello World.",
	})
}

// User demo
type User struct {
	userid   string
	login    string
	userrole string
}

//Auth JWT
func Auth() *jwt.GinJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 3,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"user_id":  v.userid,
					"login":    v.login,
					"userrole": v.userrole,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				userid:   claims["user_id"].(string),
				login:    claims["login"].(string),
				userrole: claims["userrole"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			location := loginVals.Location

			dbConnect := config.Connect()
			defer dbConnect.Close()

			loginpass := ""

			switch location {
			case "admin":
				loginpass = fmt.Sprintf("SELECT user_id, login, userrole FROM public.tuser where lower(login) = lower('%s') AND pass = '%s' and del in (0) and userrole in (1,2);", userID, password)

			case "portal":
				loginpass = fmt.Sprintf("SELECT user_id, login, userrole FROM public.tuser where lower(login) = lower('%s') AND pass = '%s' and del in (0, 2);", userID, password)

			}

			theCase := "lower"
			data, err := gosqljson.QueryDbToMap(dbConnect, theCase, loginpass)

			if len(data) == 0 {

				return nil, jwt.ErrFailedAuthentication

			}

			pool := &User{
				userid:   data[0]["user_id"],
				login:    data[0]["login"],
				userrole: data[0]["userrole"],
			}

			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("DB login auth: %s", err.Error()))
			}

			if pool.userid != "" {
				return pool, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(pool interface{}, c *gin.Context) bool {
			if v, ok := pool.(*User); ok {
				_ = v
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}

//Token get pool
func Token(c *gin.Context) {

	claims := jwt.ExtractClaims(c)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   claims,
	})
}

//Logout logout
func Logout(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	token, _ := c.Get("JWT_TOKEN")

	inserttoken := fmt.Sprintf(`INSERT INTO public.logout
	("token")
	VALUES('%s');`, token)

	_, err := dbConnect.Exec(inserttoken)

	if err != nil {
		log.Fatal("Insert token:" + err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   token,
	})

}

//Blacklist check token in blacklist
func Blacklist(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	token := jwt.GetToken(c)

	todo := fmt.Sprintf(`SELECT "token"
	FROM public.logout where token = '%s';`, token)

	var blacktoken string

	sql := dbConnect.QueryRow(todo)
	sql.Scan(&blacktoken)

	if blacktoken != "" {
		c.AbortWithStatusJSON(401, gin.H{"Error": "Your token is blacklisted"})
		return
	}
	c.Next()

}

//GetTokenInfo Get Token Info
func GetTokenInfo(c *gin.Context) {

	//data := jwt.GetToken(c)
	data := jwt.ExtractClaims(c)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

	return data

}
