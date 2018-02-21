package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	//"net/url"
	//"os"

	_ "github.com/mattn/go-sqlite3"
	"fmt"
//	"path/filepath"
//	"strings"
)

type list struct {
	Id int
	Name string
	Info string
}

type listOflist struct {
	Id int
	ListID int
	Name string
	Source string
	Date string
	ListName string
}


func index (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	rows, _ := db.Query("SELECT Id, Name, Info FROM list")
	var ls list
	var tmp []list
	for rows.Next() {
		rows.Scan(&ls.Id, &ls.Name, &ls.Info)
		tmp = append(tmp, ls)
	}
	defer db.Close()
	t := template.Must(template.ParseFiles("./templates/list.html"))
	t.ExecuteTemplate(w, "list",&tmp)
}

func edit (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	name := r.URL.Query().Get("name")
	//fmt.Println(id, name)
	rows, _ := db.Query("SELECT * FROM listoflist WHERE ListID=?",id)
	var ls listOflist
	var tmp []listOflist
	for rows.Next() {
		rows.Scan(&ls.ListID, &ls.Name, &ls.Source, &ls.Date, &ls.Id, &ls.ListName)
		tmp = append(tmp, ls)
	}
	if len(tmp) != 0 {
		var v = make(map[string]interface{})
		v["Lists"] = tmp
		v["Id"]=id
		v["Name"]=name
		//fmt.Println(v)
		defer db.Close()
		t := template.Must(template.ParseFiles("./templates/edit.html"))
		t.ExecuteTemplate(w, "edit",v)
	} else {
		defer db.Close()
		http.Redirect(w, r, "/new?id="+r.URL.Query().Get("id"), 301)
	}
}

func save1 (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	if r.Method == "POST" {
		name := r.FormValue("name")
		info := r.FormValue("info")
        insForm, err := db.Prepare("INSERT INTO list(name,info) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name,info)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
	
}

func save2 (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	name := r.URL.Query().Get("name")
	var ls listOflist
	if r.Method == "POST" {
		ls.ListID = id
		ls.Name = r.FormValue("name")
		ls.Source = r.FormValue("source")
		ls.Date = r.FormValue("date")
		ls.ListName = name
        insForm, err := db.Prepare("INSERT INTO listoflist(ListID,Name,Source,Date,ListName) VALUES(?,?,?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(ls.ListID, ls.Name, ls.Source, ls.Date, ls.ListName)
    }
	defer db.Close()
    http.Redirect(w, r, "/edit?id="+strconv.Itoa(ls.ListID)+"&name="+name, 301)
	
}

func new (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	rows, _ := db.Query("SELECT * FROM list WHERE Id=?",id)
	var ls list
	var tmp []list
	for rows.Next() {
		rows.Scan(&ls.Id, &ls.Name, &ls.Info)
		tmp = append(tmp, ls)
	}
	defer db.Close()
	t := template.Must(template.ParseFiles("./templates/new.html"))
	t.ExecuteTemplate(w, "new", tmp)
}

func update (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var ls listOflist
    if r.Method == "POST" {
		ls.Name = r.FormValue("name")
		ls.Source = r.FormValue("source")
		ls.Date = r.FormValue("date")
		ls.ListName = r.FormValue("gpname")
        insForm, err := db.Prepare("UPDATE listoflist SET Name=?,Source=?,Date=? WHERE Id=?")
        if err != nil {
            panic(err.Error())
        }
		insForm.Exec(ls.Name, ls.Source, ls.Date, id)
    }
	defer db.Close()
	
    http.Redirect(w, r, "/edit?id="+r.URL.Query().Get("listid")+"&name="+ls.ListName, 301)
}

func delete (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	//listid, _ := strconv.Atoi(r.URL.Query().Get("listid"))
	insForm, err := db.Prepare("DELETE FROM listoflist WHERE Id=?")
    if err != nil {
            panic(err.Error())
        }
    insForm.Exec(id)
    defer db.Close()
    http.Redirect(w, r, "/edit?id="+r.URL.Query().Get("listid"), 301)
}

func deletegp (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	//listid, _ := strconv.Atoi(r.URL.Query().Get("listid"))
	insForm1, _ := db.Prepare("DELETE FROM listoflist WHERE ListID=?")
	insForm2, _ := db.Prepare("DELETE FROM list WHERE Id=?")
	insForm1.Exec(id)
	insForm2.Exec(id)
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func copy (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	name := r.URL.Query().Get("name")
	rows, _ := db.Query("SELECT * FROM listoflist WHERE Id=?",id)
	var ls listOflist
	var tmp []listOflist
	for rows.Next() {
		rows.Scan(&ls.ListID, &ls.Name, &ls.Source, &ls.Date, &ls.Id, &ls.ListName)
		tmp = append(tmp, ls)
	}
	insForm, err := db.Prepare("INSERT INTO listoflist(ListID,Name,Source,Date,ListName) VALUES(?,?,?,?,?)")
    if err != nil {
        panic(err.Error())
    }
	insForm.Exec(tmp[0].ListID, tmp[0].Name, tmp[0].Source, tmp[0].Date, tmp[0].ListName)
	defer db.Close()
    http.Redirect(w, r, "/edit?id="+strconv.Itoa(ls.ListID)+"&name="+name, 301)
	
}

func export (w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./proba.db")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	rows, _ := db.Query("SELECT * FROM listoflist WHERE ListID=?",id)
	var ls listOflist
	var tmp []listOflist
	for rows.Next() {
		rows.Scan(&ls.ListID, &ls.Name, &ls.Source, &ls.Date, &ls.Id, &ls.ListName)
		tmp = append(tmp, ls)
	}
	s := ""
	for i := 0; i < len(tmp); i++ {
		s += tmp[i].Name + "\n" + tmp[i].Source + "\n" + "\n"
	}
	fmt.Println(s,id)
	w.Header().Set("Content-Disposition", "attachment; filename="+tmp[0].ListName+".log")
	w.Header().Set("Content-Type","text/plain")
	w.Write([]byte(s))
	//http.Redirect(w, r, "/", 301)
}

func main () {
	http.HandleFunc("/", index)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/save1", save1)
	http.HandleFunc("/save2", save2)
	http.HandleFunc("/new", new)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/deletegp", deletegp)
	http.HandleFunc("/copy", copy)
	http.HandleFunc("/export", export)

	http.ListenAndServe(":8081", nil)
}