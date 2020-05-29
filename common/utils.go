// Common tools and helper functions
package common

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/julienschmidt/httprouter"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// A helper function to generate random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Keep this two config private, it should not expose to open source
const NBSecretPassword = "A String Very Very Very Strong!!@##$!@#$"
const NBRandomPassword = "A String Very Very Very Niubilty!!@##$!@#4"

// A Util function to generate jwt_token which can be used in the request header
func GenToken(id uint, rule string) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		"id":   id,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"rule": rule,
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte(NBSecretPassword))
	return token
}

// My own Error type that will help return my customized Error info
//  {"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

// To handle the error returned by c.Bind in gin framework
func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		if v.Param() != "" {
			res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		} else {
			res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		}

	}
	return res
}

// Warp the error info in a object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

// Changed the c.MustBindWith() ->  c.ShouldBindWith().
// I don't want to auto return 400 when error happened.
// origin function is here: https://github.com/gin-gonic/gin/blob/master/context.go
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

// Debugging Tool / Pretty Print of Struct
func PrettyPrint(input interface{}) {
	b, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

// function for retrying task execution, using this function
// you can easily retry task over failures with specified duration.
func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}

// retry end

// Convert Http Router to Gin Routing
func ConverHttprouterToGin(f httprouter.Handle) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params httprouter.Params
		_len := len(c.Params)
		if _len == 0 {
			params = nil
		} else {
			params = ((*[1 << 10]httprouter.Param)(unsafe.Pointer(&c.Params[0])))[:_len]
		}

		f(c.Writer, c.Request, params)
	}
}

// Read Integer From String
func ReadInt(r *http.Request, param string, v int64) (int64, error) {
	p := r.FormValue(param)
	if p == "" {
		return v, nil
	}

	return strconv.ParseInt(p, 10, 64)
}

// Write Json
func WriteJSON(w http.ResponseWriter, v interface{}) {
	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(data)
}

// Read Json
func ReadJSON(r *http.Request, v interface{}) error {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, v)
}

// Common API models for using across different APIs

type MetaData struct {
	Filters    []Filter   `json:"filters"`
	Pagination Pagination `json:"pagination"`
	Order      Order      `json:"order"`
}

type Filter struct {
	Field    string `json:"field,required"`
	Value    string `json:"value,required"`
	Operator string `json:"operator,required"`
}

type Order struct {
	OrderBy   string `json:"orderBy,omitempty"`
	OrderType string `json:"orderType,omitempty"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

// dataTables Request Types
type GetAllRequest struct {
	Filters []Filter `json:"filters,omitempty"`
	Order
	Limit  int    `json:"limit,required"`
	Offset int    `json:"offset,required"`
	Query  string `json:"query,omitempty"`
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cache-Control")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cache-Control")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

// function to check if key exist in filter types
// return index if exists and return -1 if not found
func IsExistInFilters(filters []Filter, key string) int {
	for i, filter := range filters {
		if filter.Field == key {
			return i
		}
	}
	return -1
}

// by using this function we can build search queries by filter types
func BuildSearchByTags(model interface{}, filters []Filter, db *gorm.DB) *gorm.DB {
	t := reflect.TypeOf(model)
	// Iterate over all available fields and read the tag value

	// first let see if it's search request
	searchReq := -1
	for i, filter := range filters {
		if filter.Field == "allKeys" {
			searchReq = i
		}
	}

	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("gorm")
		col := ""
		if tag != "" {
			_, _ = fmt.Sscanf(tag, "column:%s", &col)
			if col != "" {
				if strings.ContainsAny(col, ";") {
					col = col[0:strings.Index(col, ";")]
				}
			}
		}
		isSearchable := field.Tag.Get("searchable")

		if searchReq != -1 {
			if isSearchable != "false" {
				if isSearchable == "string" {
					db = db.Or(col + " LIKE '%" + filters[searchReq].Value + "%'")
				}
				if isSearchable == "int" {
					db = db.Or(col + " = '" + filters[searchReq].Value + "'")
				}

			}
		} else {
			isFound := IsExistInFilters(filters, col)
			if isFound != -1 && isSearchable != "false" {
				f := filters[isFound]
				if f.Operator != "LIKE" {
					db = db.Or(f.Field + " " + f.Operator + " '" + f.Value + "'")
				} else {
					db = db.Or(f.Field + " LIKE '%" + f.Value + "%'")
				}
			}
		}

	}
	return db
}
