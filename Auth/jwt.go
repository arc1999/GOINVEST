package Auth

import (
	"GO-INVEST/errs"
	"GO-INVEST/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"os"
	"time"
	"net/http"
	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)
const (
	APP_KEY = "DAEMONS"
)
type User_temp struct{
	Uname string `json:"Username"`
	Upass string `json:"Password"`
}
var dnsstr1 =fmt.Sprintf("root:root@tcp(%s:3306)/Stocks?charset=utf8&parseTime=True",os.Getenv("CONTAINER_NAME"))
// TokenHandler is our handler to take a username and password and,
// if it's valid, return a token used for future requests.
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	var u1 User_temp
	var temp types.Users
	_ = json.NewDecoder(r.Body).Decode(&u1)
	fmt.Println(u1)
	db, err := gorm.Open("mysql", dnsstr1)
	defer db.Close()
	Db:= db.Table("users").Where("u_name = ?", u1.Uname).Find(&temp)
	if Db.RowsAffected == 0{
		err=errors.New("Incorrect Credentials")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}
	if  u1.Upass != temp.U_pass {
		w.WriteHeader(http.StatusUnauthorized)
		err=errors.New("Incorrect Credentials")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}

	// We are happy with the credentials, so build a token. We've given it
	// an expiry of 1 hour.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": temp.U_id,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

// AuthMiddleware is our middleware to check our token is valid. Returning
// a 401 status to the client if it is not valid.
func AuthMiddleware(next http.Handler) http.Handler {
	if len(APP_KEY) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(APP_KEY), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}
func Signup(w http.ResponseWriter, r *http.Request){
	var u1 types.Users
	var u2 types.Users
	_ = json.NewDecoder(r.Body).Decode(&u1)
	fmt.Println(u1)
	db, err := gorm.Open("mysql", dnsstr1)
	defer db.Close()
	temp:=db.Table("users").Where("u_name = ?",u1.U_name).Find(&u2)
	if(temp.RowsAffected!=0){
		err=errors.New("Username Already exsist ")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}
	if(err!=nil){
		err=errors.New("Database error")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}
	u1.Amount=50000.00
	db.Create(&u1)
	b, _ := json.Marshal(u1)
	_,_=fmt.Fprintf(w,"%s", b)
}