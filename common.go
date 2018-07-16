package util

import (
	"reflect"
	"fmt"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin"
)

func InArray(val interface{}, array interface{}) (exists bool) {
	exists = false

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return
			}
		}
	}

	return
}

func Print(data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Printf("%s\n", jsonData)
	return;
}

// This function will be used as replacement of JS Array.map
func Map(array []bson.M, cb func(bson.M) bson.M) []bson.M {
	for key, item := range array {
		array[key] = cb(item)
	}
	return array
}

func GetParams(ctx *gin.Context) map[string]interface{} {
	params := ctx.Query("params")
	var data interface{}
	if params != "" {
		in := []byte(params)
		json.Unmarshal(in, &data)
	}
	param, _ := data.(map[string]interface{})
	return param
}
