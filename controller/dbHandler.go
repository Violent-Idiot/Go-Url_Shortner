package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// type UrlPair struct {
// 	// Id   []byte
// 	// Name string
// 	Url string
// }

func FetchUrl(db *sql.DB, name string) string {
	query := fmt.Sprintf(`SELECT url FROM url WHERE name = "%s"`, name)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// var urls []UrlPair

	var url string
	for rows.Next() {

		err := rows.Scan(&url)
		if err != nil {
			log.Fatal(err)
		}
		// urls = append(urls, url)
	}
	// log.Print(urls)
	return url
}

func SaveUrl(db *sql.DB, Name, Url, Date string) error {

	if ifExist := checkData(db, Name); ifExist {
		return errors.New("data already exist")
	}
	// if ifDuplicateExist := checkDuplicateData(db, Name, Url); ifDuplicateExist {
	// 	return nil
	// }
	uuid, _ := uuid.NewUUID()
	// log.Print(uuidBytes, uuid)
	log.Printf("%s %s", Name, Url)
	// time := time.Date(Date)
	time, err := time.Parse(time.RFC3339, Date)
	if err != nil {
		log.Fatal(err)
	}
	monthName, day := time.Month(), time.Day()
	monthMapping := map[string]int{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}
	month := monthMapping[monthName.String()]

	query := fmt.Sprintf(`INSERT INTO url ( id, name, url) VALUES(UNHEX(REPLACE("%v", "-","")), "%s", "%s");`, uuid, Name, Url)
	db.Query(query)
	c := cron.New()
	// log.Printf("%v %v %v", month, day, time.Second())

	c.AddFunc(fmt.Sprintf("* * %v %v *", day, month), func() {
		log.Println("Cron Delete")
		deleteQuery := fmt.Sprintf(`DELETE FROM url WHERE id = UNHEX(REPLACE("%v", "-",""))`, uuid)
		_, err := db.Query(deleteQuery)
		if err != nil {
			log.Fatal(err)
		}
	})
	// inspect(c.Entries())
	c.Start()
	// defer c.Stop()
	// if err != nil {
	// 	return err
	// }
	return nil
}

func checkData(db *sql.DB, name string) bool {
	query := fmt.Sprintf(`SELECT * FROM url WHERE name = "%s"`, name)
	rows, err := db.Query(query)
	// log.Print(rows.Next())
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		return true
	}
	return false
}

func checkDuplicateData(db *sql.DB, name, url string) bool {
	query := fmt.Sprintf(`SELECT * FROM url WHERE name = "%s" AND url = "%s"`, name, url)
	rows, err := db.Query(query)
	// log.Print(rows.Next())
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		return true
	}
	return false
}
