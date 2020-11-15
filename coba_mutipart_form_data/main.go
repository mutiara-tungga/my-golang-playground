package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/healthz", Healthz)
	router.POST("/test/:id", Test)

	log.Fatal(http.ListenAndServe(":8001", router))
}

func Healthz(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ok")
}

func Test(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	rawSQLFile, rawSQLHeader, err := r.FormFile("raw_sql")
	fmt.Println(rawSQLFile == nil)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	name := r.FormValue("name")
	id := params.ByName("id")
	defer rawSQLFile.Close()

	rawSQL, err := ioutil.ReadAll(rawSQLFile)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	rawSQLString := strings.ReplaceAll(string(rawSQL), "\n", " ")

	fmt.Println("raw sqllll ", rawSQLString)
	fmt.Println("raw sql header", rawSQLHeader)
	fmt.Println("name", name)
	fmt.Println("id", id)

	res := map[string]interface{}{
		"raw_sql":        rawSQLString,
		"raw_sql_header": rawSQLHeader.Header,
		"name":           name,
		"id":             id,
		"mime_type":      http.DetectContentType(rawSQL),
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resJSON)
}
