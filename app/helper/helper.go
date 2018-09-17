package helper

import (
	"log"
	"net/http"
	"reflect"
)

//InArray func
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

//SetHeader func
func SetHeader(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
	// a := string(b)
	// log.Println(a)
	// log.Print(b)
	log.Print("Handle is ok")
}
