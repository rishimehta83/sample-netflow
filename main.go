package main

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type NetFlow struct {
	SrcApp     string `json:"src_app"`
	DestApp    string `json:"dest_app"`
	VpcID      string `json:"vpc_id"`
	BytesTx    int    `json:"bytes_tx"`
	BytesRx    int    `json:"bytes_rx"`
	Hour       int    `json:"hour"`
}


const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " +
				"password=%s dbname=%s sslmode=disable",
				host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	p, err := strconv.Atoi(r.URL.Query().Get("hour"))       

	userSql := "select distinct src_app, dest_app, vpc_id, sum(bytes_tx) " +
		   "as bytes_tx, sum(bytes_rx) as bytes_rx, hour from netflow " +
		   "where hour = $1 group by (src_app, dest_app, vpc_id, hour);"
	rows, err := db.Query(userSql, p)
	if err != nil {
		log.Fatal(err)
	}

	flows := make([]NetFlow, 0)

	for rows.Next() {
		var flow NetFlow
		err := rows.Scan(&flow.SrcApp, &flow.DestApp, &flow.VpcID, 
				 &flow.BytesTx, &flow.BytesRx, &flow.Hour)
		if err != nil {
			log.Fatal(err)
                }
		flows = append(flows, flow)
	}

	flowBytes, _ := json.MarshalIndent(flows, "", "\t")
        fmt.Println(string(flowBytes))

	w.Header().Set("Content-Type", "application/json")
	w.Write(flowBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p []NetFlow
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
       
	for _, flow := range p {
		sqlStatement := `INSERT INTO netflow (src_app, dest_app, vpc_id,
				 bytes_tx, bytes_rx, hour) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err2 := db.Exec(sqlStatement, flow.SrcApp, flow.DestApp, 
				   flow.VpcID, flow.BytesTx, flow.BytesRx, flow.Hour)

		if err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
 	   	}
        }

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}


func ReqHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case "GET":
			GETHandler(w, r)
		case "POST":
			POSTHandler(w, r)

	}
}

func main() {
	http.HandleFunc("/flows", ReqHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
