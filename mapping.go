package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// Name of the struct tag used in examples
const tagName = "validate"

// User ....
type User struct {
	ID       int    `kvdata:"id"`
	Name     string `kvdata:"name"`
	Email    string `kvdata:"email"`
	Eligible string `kvdata:"eligible"`
}

// BindMap ...
func BindMap(data map[string]interface{}, a interface{}) {
	v := reflect.ValueOf(a).Elem()
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		kvDataTag := v.Type().Field(j).Tag.Get("kvdata")
		switch v.Type().Field(j).Type.Name() {
		case "int":
			value, err := strconv.ParseInt(fmt.Sprint(data[kvDataTag]), 0, 64)
			if err != nil {
				log.Fatal(err)
			}
			f.SetInt(value)
		case "string":
			f.SetString(fmt.Sprint(data[kvDataTag]))
		case "bool":
			value, err := strconv.ParseBool(fmt.Sprint(data[kvDataTag]))
			if err != nil {
				log.Fatal(err)
			}
			f.SetBool(value)
		default:
			log.Fatal("not supported")
		}
	}
}

func main() {
	userMap := make(map[string]interface{})
	userMap["id"] = "1"
	userMap["name"] = "name"
	userMap["email"] = "email"
	userMap["eligible"] = "true"

	var user User
	BindMap(userMap, &user)
	fmt.Printf("\n%+v\n", user)
}
