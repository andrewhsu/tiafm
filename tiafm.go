package tiafm

import (
	"regexp"
	"strings"
	"unicode"
)

// tree.json
type Event struct {
	Code       string
	Name       string
	Start_date string
	End_date   string
	Track      string
}

func (e *Event) Clean() {
	e.Code = strings.TrimSuffix(e.Code, "/afm")
}

type Season struct {
	Id     string
	Year   int
	Events []Event
}

func (s *Season) Clean() {
	for i := range s.Events {
		s.Events[i].Clean()
	}
}

type Results_Tree struct {
	Seasons []Season
}

func (rt *Results_Tree) Clean() {
	for i := range rt.Seasons {
		rt.Seasons[i].Clean()
	}
}

// races-20150530.json
type Session_Race struct {
	Actual_start     string
	Actual_finish    string
	ClassID          string
	ClassName        string
	Name             string
	Raid             string
	ResultStatus     int
	Run_length       string
	Run_length_str   string
	Sanction_id      string
	Sched_length     string
	Sched_length_str string
	Scheduled_start  string
	Scheduled_finish string
	Shortcode        string
	Status           int
	Type             string
	Wave             int
}

func (sr *Session_Race) Clean() {
	r := regexp.MustCompile(` +`)
	sr.ClassName = r.ReplaceAllString(sr.ClassName, " ")
}

type Session struct {
	Start       string
	Combo_title string
	Races       []Session_Race
}

func (s *Session) Clean() {
	for i := range s.Races {
		s.Races[i].Clean()
	}
}

type Day struct {
	Date     string
	Sessions []Session
}

func (d *Day) Clean() {
	for i := range d.Sessions {
		d.Sessions[i].Clean()
	}
}

type Race_Races struct {
	Days []Day
}

func (rr *Race_Races) Clean() {
	for i := range rr.Days {
		rr.Days[i].Clean()
	}
}

// results-2900000000000071.json
type Result struct {
	BestLap            string
	BestLapOnLap       int
	ChampionshipPoints string
	Gap                string
	LeaderGap          string
	Points             string
	Position           string
	RacerID            string
	RacerName          string
	Raid               string
	Rrid               string
	SanctionID         string
	SanctionStatus     string
	Sponsors           string
	Tiid               string
	TotalLaps          interface{}
	TotalTime          string
	VehicleID          string
	VehicleName        string
}

func (r *Result) Clean() {
	r.BestLap = strings.TrimLeft(r.BestLap, " 0:")

	n := strings.Replace(r.RacerName, "*", "", -1)
	n = strings.Trim(n, " ")
	if len(n) > 0 && unicode.IsLower(rune(n[0])) {
		n = strings.Title(n)
	}

	if len(n) > 0 && !unicode.IsLetter(rune(n[0])) {
		n = ""
	}
	r.RacerName = n

	reg := regexp.MustCompile(` +`)
	r.VehicleName = reg.ReplaceAllString(r.VehicleName, " ")
	r.VehicleName = strings.Trim(r.VehicleName, " ")
}

type Race_Result struct {
	Raid    string
	Ra_type string
	Results []Result
}

func (rr *Race_Result) Clean() {
	rr.Raid = strings.ToLower(rr.Raid)
	rr.Ra_type = strings.ToLower(rr.Ra_type)
	for i := range rr.Results {
		rr.Results[i].Clean()
	}
}
