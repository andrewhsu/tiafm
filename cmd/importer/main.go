package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/andrewhsu/tiafm"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 4 {
		fmt.Println("usage: importer <schema_file> <db_file> <path_to_json_files>")
		os.Exit(1)
	}

	schema := os.Args[1]
	dbfile := os.Args[2]
	path := os.Args[3]

	tree, err := ioutil.ReadFile(path + "tree.json")
	if err != nil {
		log.Fatal(err)
	}

	var rt tiafm.Results_Tree
	err = json.Unmarshal(tree, &rt)
	if err != nil {
		log.Fatal(err)
	}
	rt.Clean()

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql, err := ioutil.ReadFile(schema)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range rt.Seasons {
		fmt.Printf("%d %s\n", s.Year, s.Id)
		for _, e := range s.Events {
			fmt.Printf("  %s %s\n", e.Code, e.Name)

			stmt, err := tx.Prepare("INSERT OR IGNORE INTO race(" +
				"Season_Year,Season_Id," +
				"Event_Code,Event_Name,Event_Start_date,Event_End_date,Event_Track," +
				"Session_Combo_title,Session_Start," +
				"Session_Race_Actual_start,Session_Race_Actual_finish," +
				"Session_Race_ClassID,Session_Race_ClassName," +
				"Session_Race_Name,Session_Race_Raid,Session_Race_ResultStatus," +
				"Session_Race_Run_length,Session_Race_Run_length_str," +
				"Session_Race_Sanction_id,Session_Race_Scheduled_start,Session_Race_Scheduled_finish,Session_Race_Shortcode," +
				"Session_Race_Status,Session_Race_Type,Session_Race_Wave" +
				") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			if err != nil {
				log.Fatal(err)
			}

			filename := fmt.Sprintf("races-%s.json", e.Code)
			log.Println(filename)
			races, err := ioutil.ReadFile(path + filename)
			if err != nil {
				log.Fatal(err)
			}

			var rr tiafm.Race_Races
			err = json.Unmarshal(races, &rr)
			if err != nil {
				log.Fatal(err)
			}
			rr.Clean()

			for _, d := range rr.Days {
				for _, ses := range d.Sessions {
					for _, sr := range ses.Races {
						_, err := stmt.Exec(s.Year, s.Id,
							e.Code, e.Name, e.Start_date, e.End_date, e.Track,
							ses.Combo_title, ses.Start,
							sr.Actual_start, sr.Actual_finish,
							sr.ClassID, sr.ClassName,
							sr.Name, sr.Raid, sr.ResultStatus,
							sr.Run_length, sr.Run_length_str,
							sr.Sanction_id, sr.Scheduled_start, sr.Scheduled_finish, sr.Shortcode,
							sr.Status, strings.ToLower(sr.Type), sr.Wave)
						if err != nil {
							log.Fatal(err)
						}

						filename := fmt.Sprintf("results-%s.json", sr.Raid)
						log.Println(filename)
						race_result, err := ioutil.ReadFile(path + filename)
						if err != nil {
							log.Fatal(err)
						}

						var rr tiafm.Race_Result
						err = json.Unmarshal(race_result, &rr)
						if err != nil {
							log.Fatal(err)
						}
						rr.Clean()

						stmt2, err := tx.Prepare("INSERT OR IGNORE INTO result(" +
							"Race_Result_Ra_type,Result_BestLap,Result_BestLapOnLap," +
							"Result_ChampionshipPoints,Result_Gap,Result_LeaderGap," +
							"Result_Points,Result_Position," +
							"Result_RacerID,Result_RacerName," +
							"Result_Raid,Result_Rrid," +
							"Result_SanctionID,Result_SanctionStatus," +
							"Result_Sponsors,Result_Tiid," +
							"Result_TotalLaps,Result_TotalTime," +
							"Result_VehicleID,Result_VehicleName" +
							") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
						if err != nil {
							log.Fatal(err)
						}

						for _, res := range rr.Results {
							_, err := stmt2.Exec(rr.Ra_type, res.BestLap, res.BestLapOnLap,
								res.ChampionshipPoints, res.Gap, res.LeaderGap,
								res.Points, res.Position,
								res.RacerID, res.RacerName,
								res.Raid, res.Rrid,
								res.SanctionID, res.SanctionStatus,
								res.Sponsors, res.Tiid,
								res.TotalLaps, res.TotalTime,
								res.VehicleID, res.VehicleName)
							if err != nil {
								log.Fatal(err)
							}
						}

						stmt2.Close()

					}
				}
			}
			stmt.Close()

		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
