package src

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"os"
)

const (
	PriceTable = "pricetable"
)

type DataBaseWork struct{
	DBConnection *sql.DB
	Config       *Config
}

func (db *DataBaseWork) addRecord(b *Body) error {
	log.Println("In addRecord")
	row := db.DBConnection.QueryRow("SELECT Mail,link,price FROM PriceTable WHERE link = $1", b.Link)
	var tmp1, tmp2, tmp3 string
	err := row.Scan(&tmp1, &tmp2, &tmp3)
	runGoroutine := false
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
			return err
		}
		// Record doesn't exist
		runGoroutine = true
	}

	curPrice, err := getPrice(&b.Link)
	_, err = db.DBConnection.Exec("INSERT INTO PriceTable (Mail, link, price) VALUES ($1, $2, $3)", b.Mail, b.Link, curPrice)
	if err != nil {
		log.Println(err)
		return err
	}
	if runGoroutine {
		go UpdatingLinkProcess(db, &b.Link, curPrice)
	}
	return nil
}

func SetDB() (*DataBaseWork, error) {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//_, err = db.Exec("CREATE TABLE pricetable (Mail TEXT, link TEXT, price VARCHAR(32), unique(Mail, link))")
	//if err != nil {
	//	log.Println(err)
	//}
	return &DataBaseWork{
		DBConnection: db,
		Config:       new(Config),
	}, nil
}
