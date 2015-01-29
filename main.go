package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fcoury/gridhook/models"
	"github.com/go-martini/martini"

	_ "github.com/go-sql-driver/mysql"
)

func SetupDB() *sql.DB {
	db, err := sql.Open("mysql", "root@/crm_bliss?parseTime=true")
	PanicIf(err)

	return db
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(500)
	fmt.Fprintf(rw, "Server Error: %s\n", err)
}

func HandleSuccess(rw http.ResponseWriter) {
	rw.WriteHeader(200)
}

var port int

func main() {
	flag.IntVar(&port, "port", 5000, "The port to run the SendGrid Webhook listener on")
	flag.Parse()

	m := martini.Classic()
	m.Map(SetupDB())

	m.Post("/sendgrid/event", func(db *sql.DB, r *http.Request, rw http.ResponseWriter) {
		var m map[string]interface{}

		body, err := ioutil.ReadAll(r.Body)
		PanicIf(err)

		err = json.Unmarshal(body, &m)
		if err != nil {
			rw.WriteHeader(400)
			fmt.Fprintf(rw, "Bad request: %s\n", err)
			return
		}

		if val, ok := m["email_event_id"]; ok {
			uniqueId := val.(string)
			e, err := models.FindEmailEventByUniqueId(db, uniqueId)
			if err != nil {
				HandleError(rw, err)
				return
			}

			if e == nil {
				fmt.Fprintf(rw, "EmailEvent not found: %s\n", uniqueId)
				return
			}

			e.FetchStatuses(db)
			last, err := e.FindMostRecentStatus(db)
			PanicIf(err)
			fmt.Printf("Status: %s\n", last.Status)

			if m["timestamp"] == nil {
				HandleSuccess(rw)
				return
			}

			timestamp := time.Unix(int64(m["timestamp"].(float64)), 0)

			status := m["event"].(string)

			err = e.InsertStatus(db, timestamp, status)
			PanicIf(err)

			if last.Timestamp.After(timestamp) {
				return
			}

			e.Status = status
			err = e.Update(db)
			PanicIf(err)

			fmt.Fprintf(rw, "New status: %s\n", status)
		} else {
			fmt.Fprintf(rw, "200 OK")
		}
	})

	host := fmt.Sprintf("localhost:%v", port)
	m.RunOnAddr(host)
}
