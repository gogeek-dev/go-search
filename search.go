package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func dbConn() (db *sql.DB) {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er.Error())
	}
	dbDriver := os.Getenv("DB_Driver")
	dbUser := os.Getenv("DB_User")
	dbPass := os.Getenv("DB_Password")
	dbName := os.Getenv("DB_Name")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

type person struct {
	ID   int
	Name string
	City string
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM person")
	if err != nil {
		panic(err.Error())
	}
	emp := person{}
	res := []person{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.ID = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}

	tmpl.ExecuteTemplate(w, "index.html", res)
	defer db.Close()

}
func search(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if r.Method == "POST" {
		keyword := r.FormValue("search")

		selDB, err := db.Query("SELECT * FROM person where name  like  ? or city like ?", "%"+keyword+"%", "%"+keyword+"%")
		if err != nil {
			panic(err.Error())
		}
		emp := person{}
		res := []person{}
		for selDB.Next() {
			var id int
			var name, city string
			err = selDB.Scan(&id, &name, &city)
			if err != nil {
				panic(err.Error())
			}
			emp.ID = id
			emp.Name = name
			emp.City = city
			res = append(res, emp)
		}

		tmpl.ExecuteTemplate(w, "index.html", res)
		defer db.Close()
	}
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/searchitem", search)
	http.ListenAndServe(":8080", nil)
}
