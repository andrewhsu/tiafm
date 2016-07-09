package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func racerHandler(w http.ResponseWriter, r *http.Request) {
	tiid := r.URL.Path[len("/r/"):]
	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type Result struct {
		Event   string
		Race    string
		Pos     string
		Num     string
		Lic     string
		Best    string
		Pts     string
		Laps    int
		Vehicle string
	}

	rows, err := db.Query("SELECT Result_RacerName,Event_Name,Session_Race_ClassName," +
		"Result_Position,Result_SanctionID,Result_SanctionStatus," +
		"Result_BestLap,Result_Points,Result_TotalLaps,Result_VehicleName " +
		"FROM result JOIN race ON Session_Race_Raid=Result_Raid " +
		"WHERE Result_Tiid=" + tiid + " AND Race_Result_Ra_type='race' " +
		"ORDER BY Session_Race_Scheduled_finish DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	context := struct {
		Name    string
		Results []Result
	}{}
	var result Result
	for rows.Next() {
		rows.Scan(&context.Name, &result.Event, &result.Race,
			&result.Pos, &result.Num, &result.Lic,
			&result.Best, &result.Pts, &result.Laps, &result.Vehicle)
		switch result.Pos {
		case "9994":
			result.Pos = "DNF"
		case "9995":
			result.Pos = "DQ"
		case "9997":
			result.Pos = "DQ"
		case "9998":
			result.Pos = "DNS"
		}
		context.Results = append(context.Results, result)
	}

	t, err := template.ParseFiles("racer.html.template")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Fatal(err)
	}
}

func classHandler(w http.ResponseWriter, r *http.Request) {
	classid := r.URL.Path[len("/c/"):]
	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type Line struct {
		Time    string
		Date    string
		Num     string
		Name    string
		Vehicle string
	}

	query := "SELECT Session_Race_ClassName,Result_BestLap,Session_Race_Scheduled_start,Result_SanctionID,Result_RacerName,Result_VehicleName " +
		"FROM race JOIN result ON Session_Race_Raid=Result_Raid " +
		"WHERE Result_BestLap!='' AND Event_Track=9 AND Session_Race_ClassID=" + classid + " " +
		"ORDER BY Result_BestLap LIMIT 10"
	log.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	context := struct {
		Class string
		Bw    []Line
		So    []Line
		Th    []Line
	}{}
	var line Line
	for rows.Next() {
		rows.Scan(&context.Class, &line.Time, &line.Date, &line.Num, &line.Name, &line.Name, &line.Vehicle)
		log.Println(line.Vehicle)
		context.Bw = append(context.Bw, line)
	}

	t, err := template.ParseFiles("class.html.template")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Fatal(err)
	}
}

func classesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type Class struct {
		Id   string
		Name string
	}

	rows, err := db.Query("SELECT DISTINCT Session_Race_ClassID,Session_Race_ClassName " +
		"FROM race " +
		"WHERE Session_Race_Type='race' " +
		"ORDER BY Session_Race_ClassName")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	context := struct{ Classes []Class }{}
	var class Class
	for rows.Next() {
		rows.Scan(&class.Id, &class.Name)
		context.Classes = append(context.Classes, class)
	}

	t, err := template.ParseFiles("classes.html.template")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Fatal(err)
	}
}

func racersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type Racer struct {
		Id     uint64
		Name   string
		Number string
	}

	rows, err := db.Query("SELECT Result_Tiid,Result_RacerName,group_concat(DISTINCT Result_SanctionID)" +
		"FROM result WHERE Result_RacerName != '' AND Result_SanctionID != '' " +
		"GROUP BY Result_RacerName " +
		"ORDER BY Result_RacerName")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	context := struct{ Racers []Racer }{}
	var racer Racer
	for rows.Next() {
		rows.Scan(&racer.Id, &racer.Name, &racer.Number)
		racer.Number = strings.Replace(racer.Number, ",", " #", -1)
		context.Racers = append(context.Racers, racer)
	}

	t, err := template.ParseFiles("racers.html.template")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/c", classesHandler)
	http.HandleFunc("/c/", classHandler)
	http.HandleFunc("/r", racersHandler)
	http.HandleFunc("/r/", racerHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
