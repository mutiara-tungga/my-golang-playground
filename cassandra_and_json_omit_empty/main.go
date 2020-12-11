package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocql/gocql"

	"github.com/julienschmidt/httprouter"
)

type ReqBody struct {
	Schedule ReqRuleSchedule `json:"rule_schedule"`
	ID       string          `json:"rule_id"`
}

type ReqRuleSchedule struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Frequency string `json:"frequency"`
	Interval  int    `json:"interval"`
	WeekDays  []int  `json:"week_days"`
	MonthDays []int  `json:"month_days"`
	Hours     []int  `json:"hours"`
}

type EntRuleSchedule struct {
	StartDate string
	EndDate   string
	Frequency string
	Interval  int
	WeekDays  []int
	MonthDays []int
	Hours     []int
}

type RepoRuleSchedule struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Frequency string `json:"frequency,omitempty"`
	Interval  int    `json:"interval,omitempty"`
	WeekDays  []int  `json:"week_days,omitempty"`
	MonthDays []int  `json:"month_days,omitempty"`
	Hours     []int  `json:"hours,omitempty"`
}

type RespRuleSchedule struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Frequency string `json:"frequency,omitempty"`
	Interval  int    `json:"interval,omitempty"`
	WeekDays  []int  `json:"week_days,omitempty"`
	MonthDays []int  `json:"month_days,omitempty"`
	Hours     []int  `json:"hours,omitempty"`
}

type RespBody struct {
	ID       string            `json:"id"`
	Schedule *RespRuleSchedule `json:"rule_schedule,omitempty"`
}

func main() {
	router := httprouter.New()
	router.GET("/get/:rule_id/:include", Get)
	router.POST("/insert", Insert)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Cassandra() *gocql.Session {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "test2"
	cluster.Consistency = gocql.LocalOne
	session, _ := cluster.CreateSession()

	return session
}

func Insert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cassandra := Cassandra()

	var reqBody ReqBody
	err := decodeReqBody(r, &reqBody)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	entRS := EntRuleSchedule{
		StartDate: reqBody.Schedule.StartDate,
		EndDate:   reqBody.Schedule.EndDate,
		Frequency: reqBody.Schedule.Frequency,
		Interval:  reqBody.Schedule.Interval,
		WeekDays:  reqBody.Schedule.WeekDays,
		MonthDays: reqBody.Schedule.MonthDays,
		Hours:     reqBody.Schedule.Hours,
	}

	repoRS := RepoRuleSchedule{
		StartDate: entRS.StartDate,
		EndDate:   entRS.EndDate,
		Frequency: entRS.Frequency,
		Interval:  entRS.Interval,
		WeekDays:  entRS.WeekDays,
		MonthDays: entRS.MonthDays,
		Hours:     entRS.Hours,
	}

	stringRS, err := json.Marshal(repoRS)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = cassandra.Query(`
		INSERT INTO table2 (
			rule_id,
			rule_schedule
		) VALUES (?, ?)`,
		reqBody.ID,
		stringRS,
	).Exec()

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("success"))
}

func Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cassandra := Cassandra()

	ruleID := ps.ByName("rule_id")
	include := ps.ByName("include")

	var rawRS string
	q := `SELECT rule_schedule FROM table2 WHERE rule_id = ?`
	err := cassandra.Query(q, ruleID).Scan(&rawRS)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var repoRS RepoRuleSchedule
	err = json.Unmarshal([]byte(rawRS), &repoRS)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	entRS := EntRuleSchedule{
		StartDate: repoRS.StartDate,
		EndDate:   repoRS.EndDate,
		Frequency: repoRS.Frequency,
		Interval:  repoRS.Interval,
		WeekDays:  repoRS.WeekDays,
		MonthDays: repoRS.MonthDays,
		Hours:     repoRS.Hours,
	}

	respRS := RespRuleSchedule{
		StartDate: entRS.StartDate,
		EndDate:   entRS.EndDate,
		Frequency: entRS.Frequency,
		Interval:  entRS.Interval,
		WeekDays:  entRS.WeekDays,
		MonthDays: entRS.MonthDays,
		Hours:     entRS.Hours,
	}

	var respBody RespBody
	respBody.ID = ruleID
	if include == "true" {
		respBody.Schedule = &respRS
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
}

func decodeReqBody(req *http.Request, object interface{}) error {
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	return decoder.Decode(&object)
}
