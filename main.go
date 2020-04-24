package main

import (
	"GO-INVEST/Auth"
	"GO-INVEST/Request"
	"GO-INVEST/Response"
	"GO-INVEST/errs"
	"GO-INVEST/types"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"log"
"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/ziutek/mymysql/godrv"
	_ "log"
	"net/http"

	"time"
)

var db *gorm.DB
var req requests.BuyStockRequest
var err error
var temp types.Stock
var temps types.Stocks
var stk types.PStocks
var pstk types.PStock
var dnsstr1 =fmt.Sprintf("root:root@tcp(%s:3306)/Stocks?charset=utf8&parseTime=True",os.Getenv("CONTAINER_NAME"))

func main() {
	dnsstr:=fmt.Sprintf("root:root@tcp(%s:3306)/",os.Getenv("CONTAINER_NAME"))
	fmt.Println(dnsstr)
	dba, err := gorm.Open("mysql", dnsstr)
	dba.Exec("CREATE DATABASE IF NOT EXISTS"+" Stocks")
	dba.Close()

	db, err = gorm.Open("mysql", dnsstr1)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	if(!db.HasTable(&types.Stock{}) ) {
		db.AutoMigrate(&types.Stock{})
	}
	if(!db.HasTable(&types.Users{}) ) {
		db.AutoMigrate(&types.Users{})
	}


	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	cors := cors.New(cors.Options{

		AllowedOrigins:   []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Route("/login", func(r chi.Router) {
		r.Post("/",Auth.TokenHandler)
		r.Post("/signup",Auth.Signup)
	})
	r.Route("/stocks", func(r chi.Router) {
		r.Use(Auth.AuthMiddleware)
		r.Post("/", BuyStock)
		r.Get("/", ListAllStocks)
		r.Get("/balance",CheckAmount)
		r.Get("/invalue",InvestedValue)
		r.Get("/profile",Profile)
		r.Get("/market",market)
		r.Get("/prvalue",predictedstocks)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(ArticleCtx)
			r.Get("/", ListOneStock)
			r.Delete("/", DeleteStock)
		})
	})
	r.Route("/lock", func(r chi.Router) {
		r.Use(Auth.AuthMiddleware)
		r.Get("/{id}",lockstock)
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "id")

		db, err = gorm.Open("mysql", dnsstr1)
		if err != nil {
			err=errors.New("DB Connection error")
			_ = render.Render(w, r, errs.ErrRender(err))
		}else{
			fmt.Println("Connection established")
		}
		var temp types.Stock
		Db:= db.Table("stocks").Where("stkid = ?", ID).Find(&temp)

		if Db.RowsAffected == 0{
			err=errors.New("ID not Found")
			_ = render.Render(w, r, errs.ErrRender(err))
			return
		} else{
			ctx := context.WithValue(r.Context(), "stock", Db )
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
func Profile(writer http.ResponseWriter, request *http.Request) {
	var user types.Users
	token := request.Context().Value("user").(*jwt.Token)
	a:=token.Claims.(jwt.MapClaims)
	id:=int(a["user"].(float64))
	err = render.Bind(request, &req)
	db, err = gorm.Open("mysql", dnsstr1)
	Db:=db.Table("users").Where("u_id=?",id).Find(&user)

	if Db.RowsAffected == 0{
		err=errors.New("User not found")
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	fmt.Println(user)
	b, _ := json.Marshal(user)
	_,_=fmt.Fprintf(writer,"%s", b)
	defer db.Close()
}
func BuyStock(writer http.ResponseWriter, request *http.Request) {
	var user types.Users
	token := request.Context().Value("user").(*jwt.Token)
	a:=token.Claims.(jwt.MapClaims)
	id:=int(a["user"].(float64))
	err = render.Bind(request, &req)
	db, err = gorm.Open("mysql", dnsstr1)
	Db:=db.Table("users").Where("u_id=?",id).Find(&user)

	if Db.RowsAffected == 0{
		err=errors.New("User not found")
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	var  req_amount  = float64(req.Initcost * float32(req.Quant))
	fmt.Println(user)
	avai_amount:= user.Amount
	if(avai_amount< req_amount){
		err=errors.New("Funds insufficent")
		_ = render.Render(writer, request, errs.ErrRender(err))
		fmt.Fprint(writer,user)
		return
	}else{
		user.Amount=user.Amount - req_amount
		Db.Save(&user)
	}
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	temp := *req.Stock
	temp.Purchased_on =time.Now().Local().Add(time.Hour * time.Duration(5) +
		time.Minute * time.Duration(30) +
		time.Second * time.Duration(0))
	temp.U_id=id
	db.Create(&temp)
	b, _ := json.Marshal(temp)
	_,_=fmt.Fprintf(writer,"%s", b)
}
func lockstock(writer http.ResponseWriter,request *http.Request){
	ID := chi.URLParam(request, "id")
	driver := mysql.MySQLDriver{}
	_ = driver
	var dbs=fmt.Sprintf("%s:3306)/major1db",os.Getenv("AWSDB"))
	dbaws, err := sql.Open("mysql", dbs)
	fmt.Println(err)
	var id int
	var cprice, sprice float32
	var name string
	err = dbaws.QueryRow("SELECT * FROM stock_predicted where stkid = ?", ID).Scan(&id, &name, &cprice, &sprice)
	if err != nil {
		fmt.Println(err)
	}
		pstk.Stk_id = id
		pstk.Stk_name = name
		pstk.Cur_price = cprice
		pstk.Pr_price = sprice

	b, _ := json.Marshal(pstk)
	stk = nil
	_, _ = fmt.Fprintf(writer, "%s", b)
	if err != nil {
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	defer dbaws.Close()
}
func CheckAmount(writer http.ResponseWriter, request *http.Request) {
	var user types.Users
	token := request.Context().Value("user").(*jwt.Token)
	a:=token.Claims.(jwt.MapClaims)
	id:=int(a["user"].(float64))
	err = render.Bind(request, &req)
	db, err = gorm.Open("mysql", dnsstr1)
	Db:=db.Table("users").Where("u_id=?",id).Find(&user)
	fmt.Fprint(writer,user.Amount)
	if Db.RowsAffected == 0{
		err=errors.New("User not found")
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
}
func ListOneStock(writer http.ResponseWriter, request *http.Request) {
	db := request.Context().Value("stock").(*gorm.DB)
	Db:= db.Find(&temp)
	if(Db.RowsAffected == 0){
		err:=errors.New("Stock not found")
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}else{
		b, _ := json.Marshal(temp)
		fmt.Println(b)
		_,_=fmt.Fprintf(writer,"%s", b)
	}
}
func predictedstocks(writer http.ResponseWriter, request *http.Request) {
	driver := mysql.MySQLDriver{}
	_ = driver
	dbaws, err := sql.Open("mysql", "admin:ayush123@tcp(major1db.c7s0raba2glo.ap-south-1.rds.amazonaws.com:3306)/major1db")
	fmt.Println(err)
	selDB, err := dbaws.Query("SELECT * FROM stock_predicted")
	if err != nil {
		fmt.Println(err)
	}

	for selDB.Next() {
		var id int
		var cprice, sprice float32
		var name string
		err = selDB.Scan(&id, &name, &cprice, &sprice)
		if err != nil {
			panic(err.Error())
		}
		pstk.Stk_id = id
		pstk.Stk_name = name
		pstk.Cur_price = cprice
		pstk.Pr_price = sprice

		stk = append(stk, pstk)
	}
	b, _ := json.Marshal(stk)
	stk = nil
	_, _ = fmt.Fprintf(writer, "%s", b)
	if err != nil {
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	defer dbaws.Close()
}
func market(writer http.ResponseWriter, request *http.Request) {
	var market1 types.Market
	var markets1 types.Markets
	driver := mysql.MySQLDriver{}
	_ = driver
	dbaws, err := sql.Open("mysql", "admin:ayush123@tcp(major1db.c7s0raba2glo.ap-south-1.rds.amazonaws.com:3306)/major1db")
	fmt.Println(err)
	selDB, err := dbaws.Query("SELECT * FROM market")
	if err != nil {
		fmt.Println(err)
	}

	for selDB.Next() {
		var id ,vol int
		var oprice, cprice, high, low float64
		var name,bid,ask string
		err = selDB.Scan(&id, &name, &oprice, &cprice,&high,&low,&bid,&ask,&vol)
		if err != nil {
			panic(err.Error())
		}
			market1.Stkid=id
			market1.Stk_name=name
			market1.Open_price=oprice
			market1.Close_price=cprice
			market1.Low=low
			market1.High=high
			market1.Bid=bid
			market1.Ask=ask
			market1.Volume=vol
		markets1 = append(markets1 , market1)
	}
	b, _ := json.Marshal(markets1)
	stk = nil
	// Convert bytes to string
	// .
	//s := string(b)
	//fmt.Println(b)
	_, _ = fmt.Fprintf(writer, "%s", b)
	if err != nil {
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	defer dbaws.Close()
}
func ListAllStocks(writer http.ResponseWriter, request *http.Request) {

	token := request.Context().Value("user").(*jwt.Token)
	a:=token.Claims.(jwt.MapClaims)
	id:=int(a["user"].(float64))
	db, err = gorm.Open("mysql", dnsstr1)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	var user types.Users
	Dba:=db.Table("users").Where("u_id=?",id).Find(&user)

	if Dba.RowsAffected == 0{
		err=errors.New("User not found")
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	fmt.Println(user)
	Db:=db.Table("stocks").Where("u_id=?",id).Find(&temps)
	if Db.RowsAffected == 0{
		err=errors.New("Stocks not found")
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	fmt.Println(temps)
	b, _ := json.Marshal(temps)
	_,_=fmt.Fprintf(writer,"%s", b)
	if err != nil {
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
}
func InvestedValue(writer http.ResponseWriter, request *http.Request) {
	token := request.Context().Value("user").(*jwt.Token)
	a:=token.Claims.(jwt.MapClaims)
	id:=int(a["user"].(float64))
	db, err = gorm.Open("mysql", dnsstr1)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	db.Table("stocks").Where("u_id=?",id).Find(&temps)
	var x float32  = 0

	for _, a := range temps {
		x += a.Initcost*float32(a.Quant)
	}
	_, _ = fmt.Fprintf(writer, "%g", x)

}
func DeleteStock(writer http.ResponseWriter, request *http.Request) {
	db := request.Context().Value("stock").(*gorm.DB)
	Db:= db.Delete(&temp)
	if(Db.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errs.ErrRender(err))
		return
	}else{
		_=render.Render(writer, request, responses.ListStock(&temp))
	}
}