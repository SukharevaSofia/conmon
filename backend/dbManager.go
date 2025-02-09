package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	log.Println("connectDB invoked")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("Could not open the database: ", err)
		return nil, err
	} else {
		log.Println("Opened: ", os.Getenv("DATABASE_URL"))
	}

	return db, nil
}

func initDB() {
	var db *sql.DB
	var err error
	for {
		log.Println("Attempting connection to DB")
		db, err = connectDB()
		if err == nil {
			break
		}
		time.Sleep(time.Microsecond * 100)
	}
	defer db.Close()

	for {
		log.Println("Table creation")
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS containers" +
			"(ip text PRIMARY KEY UNIQUE NOT NULL,ping_time timestamp, last_success timestamp)")
		if err != nil {
			log.Println("Table creation error:", err)
      continue
		} else {
			log.Println("Table conainers created")
		}
		log.Println("createDb done;")
    break;
	}
}

func readDB() ([]tableRow, error) {
	log.Println("readAll invoked")
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM containers")
	if err != nil {
		log.Println("Could read from the database: ", err)
		return nil, err
	}

	tempRow := tableRow{}
	rows := make([]tableRow, 0)
	for res.Next() {
		res.Scan(&tempRow.IP, &tempRow.Ping_time, &tempRow.Last_success)
		rows = append(rows, tempRow)
	}

	return rows, nil
}

func readDBByIP(ip string) ([]tableRow, error) {
	log.Println("readByIP invoked")
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM containers WHERE ip=$1", ip)
	if err != nil {
		log.Println("Could read from the database: ", err)
		return nil, err
	}

	tempRow := tableRow{}
	var rows []tableRow
	for res.Next() {
		res.Scan(&tempRow.IP, &tempRow.Ping_time, &tempRow.Last_success)
		rows = append(rows, tempRow)
	}

	return rows, nil
}

func addToDB(ipList []string) error {
	log.Println("addToDB invoked")

	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var last_success time.Time
	ping_time := time.Now()
	_, err = db.Exec("UPDATE containers SET ping_time = $1", ping_time)
	if err != nil {
		return err
	}

	for _, ip := range ipList {
		last_success = ping_time
		_, err = db.Exec("INSERT INTO containers VALUES ($1, $2, $3) ON CONFLICT (ip) DO UPDATE SET ping_time=$4, last_success=$5",
			ip, ping_time, last_success, ping_time, last_success)
		if err != nil {
			log.Println("Couldn't read from the database: ", err)
			return err
		}
	}
	return nil
}
