package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getTable(respWriter http.ResponseWriter, req *http.Request) {
	log.Println("getTable invoked")

	if req.Method != "GET" {
		log.Println("bad request")
		respWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	table, err := readDB()
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonedTable, err := json.Marshal(table)
	if err != nil {
		log.Println(err)
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(jsonedTable)
}

func postTableUpload(respWriter http.ResponseWriter, req *http.Request) {
	log.Println("postTable invoked")

	if req.Method != "POST" {
		log.Println("bad request")
		respWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var ipList []string
	err := decoder.Decode(&ipList)
	if err != nil {
		panic(err)
	}

  log.Println(ipList)
	if ipList == nil {
    log.Println("ipList is nil")
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
  addToDB(ipList)
	respWriter.WriteHeader(http.StatusOK)
}
