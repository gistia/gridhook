package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func main() {
	m := martini.Classic()
	m.Map(SetupDB())

	m.Post("/", func(db *sql.DB, r *http.Request, rw http.ResponseWriter) {
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

			status := m["event"]
			e.Status = status.(string)
			err = e.Update(db)
			PanicIf(err)

			fmt.Fprintf(rw, "New status: %s\n", status)
		} else {
			fmt.Fprintf(rw, "200 OK")
		}
	})

	m.Run()
}
