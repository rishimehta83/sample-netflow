package main

import (
	"database/sql"
	"encoding/json"
        "strconv"
//        "io/ioutil"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)


type Person struct {
	SrcApp     string `json:"src_app"`
	DestApp string `json:"dest_app"`
}

type NetFlow struct {
	SrcApp     string `json:"src_app"`
}

type NetFlow2 struct {
	SrcApp     string `json:"src_app"`
	DestApp    string `json:"dest_app"`
	VpcID      string `json:"vpc_id"`
        BytesTx   int  `json:"bytes_tx"`
        BytesRx   int `json:"bytes_rx"`
        Hour       int `json:"hour"`
}


const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func OpenConnection() *sql.DB {
        fmt.Println("Connection Established\n")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
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
        fmt.Printf("Rhea Mehta GET\n")

        p, err := strconv.Atoi(r.URL.Query().Get("hour"))       
        
        userSql := "SELECT * FROM netflow WHERE hour = $1" 
        rows, err := db.Query(userSql, p)
	if err != nil {
		log.Fatal(err)
	}

	var flows []NetFlow2

	for rows.Next() {
		var flow NetFlow2
		err := rows.Scan(&flow.SrcApp, &flow.DestApp, &flow.VpcID, &flow.BytesTx, &flow.BytesRx, &flow.Hour)
                if err != nil {
		   log.Fatal(err)
                }
                fmt.Printf("%+v\n", flow)
		flows = append(flows, flow)
	}

	flowBytes, _ := json.MarshalIndent(flows, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(flowBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
        fmt.Printf("Rhea Mehta POST\n")

	var p []NetFlow2
        fmt.Printf("Rhea Mehta inserting in the table\n")
	err := json.NewDecoder(r.Body).Decode(&p)
        fmt.Printf("Rhea Mehta inserting in the table\n")
	if err != nil {
                fmt.Printf("Rhea Mehta inserting in the table error\n")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO netflow (src_app, dest_app, vpc_id, bytes_tx, bytes_rx, hour) VALUES ($1, $2, $3, $4, $5, $6)`
       	_, err2 := db.Exec(sqlStatement, p[0].SrcApp, p[0].DestApp, p[0].VpcID, p[0].BytesTx, p[0].BytesRx, p[0].Hour)
        fmt.Printf("Rhea Mehta inserting in the table\n")
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

        fmt.Printf("Rhea Mehta inserting in the table done\n")
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func main() {
        fmt.Println("Received Request\n")
	http.HandleFunc("/flows", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
