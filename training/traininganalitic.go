package training

import (
	"PortalMGTNIIP/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//Getpooltrainingbyyear by YEAR
func Getpooltrainingbyyear(c *gin.Context) {

	year_text := ""
	month_text := ""
	is_external_text := ""
	end_monthf_text := ""

	year := c.Query("year")
	month := c.Query("month")
	is_external := c.Query("is_external")
	end_monthf := c.Query("end_monthf")
	end_monthl := c.Query("end_monthl")

	if year != "" {
		year_text = fmt.Sprintf(`where EXTRACT(YEAR from (training.dates_json -> 0 ->> 'date_start')::date) = %s `, year)
	}

	if month != "" {
		month_text = fmt.Sprintf(`and EXTRACT(month from (training.dates_json -> 0 ->> 'date_start')::date) = %s `, month)
	}
	if is_external != "" {
		is_external = fmt.Sprintf(`and is_external = %s `, is_external)
	}

	if end_monthf != "" && end_monthl == "" {
		end_monthf_text = fmt.Sprintf(`and EXTRACT(month from (training.dates_json -> 0 ->> 'date_end')::date) = %s `, end_monthf)
	}

	if end_monthf != "" && end_monthl != "" {

		end_monthf_text = fmt.Sprintf(`and EXTRACT(month from (training.dates_json -> 0 ->> 'date_end')::date) between %s and %s`, end_monthf, end_monthl)

	}
	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`SELECT *
	FROM public.training %s %s %s %s
	order by cast(training.dates_json -> 0 ->> 'date_start' as timestamp) desc;
	`, year_text, month_text, is_external_text, end_monthf_text)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

	return

}

func getusers(users []string) map[string][]string {

	var result = map[string][]string{}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`SELECT user_id, dep_id
	FROM public.tuser;`)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)

	}

	for i := range data {

		if searchinlist(users, data[i]["user_id"]) {

			if _, ok := result[data[i]["dep_id"]]; ok {

				result[data[i]["dep_id"]] = append(result[data[i]["dep_id"]], data[i]["user_id"])

			} else {
				result[data[i]["dep_id"]] = []string{}
				result[data[i]["dep_id"]] = append(result[data[i]["dep_id"]], data[i]["user_id"])
			}

		}
	}
	return result
}

func searchinlist(s []string, str string) bool {

	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false

}

//Users from BD
type Users struct {
	User  string `json:"user"`
	Chief string `json:"chief"`
}

func getusersintrain(year string, end_monthf string, end_monthl string) []string {

	year_text := ""
	month_text := ""

	if year != "" {
		year_text = fmt.Sprintf(` and EXTRACT(YEAR from (training.dates_json -> 0 ->> 'date_start')::date) = %s `, year)
	}

	if end_monthf != "" && end_monthl != "" {

		month_text = fmt.Sprintf(`and EXTRACT(month from (training.dates_json -> 0 ->> 'date_end')::date) between %s and %s`, end_monthf, end_monthl)

	}

	var result = []string{}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`select users 
	FROM public.training where users::text != '[]' %s %s ;`, year_text, month_text)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)

	}

	for i := range data {

		var users []*Users

		json.Unmarshal([]byte(data[i]["users"]), &users)

		for _, item := range users {

			result = append(result, item.User)

		}

	}
	log.Print(result)
	return result

}

//Getpoolusersbydep Get pool users by dep
func Getpoolusersbydep(c *gin.Context) {

	year := c.Query("year")
	end_monthf := c.Query("end_monthf")
	end_monthl := c.Query("end_monthl")

	users := getusersintrain(year, end_monthf, end_monthl)

	data := getusers(users)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

}

//TrainingTime Training Time of users
type TrainingTime struct {
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	IsOnline  string `json:"is_online"`
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
}

func getdepbyuser() map[string]string {

	var result = map[string]string{}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`select tuser.user_id, tdep.name
	FROM public.tdep tdep, public.tuser tuser
	WHERE 
		tuser.dep_id = tdep.dep_id;`)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)

	}

	for _, items := range data {

		result[items["user_id"]] = items["name"]

	}

	return result

}

func allUsersIntraining() map[string]map[int]map[string]string {

	usersindeps := getdepbyuser()

	var result = map[string]map[int]map[string]string{}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`select training.users, training.dates_json, training.is_external,
	trainingtopic.descr as topic_descr, trainingtopic.title as topic_title
	FROM public.training, public.trainingtopic where training.users::text != '[]' 
	and trainingtopic.topic_id = training.topic_id;`)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)

	}

	for i := range data {

		var users []*Users
		var dates_json []*TrainingTime

		json.Unmarshal([]byte(data[i]["users"]), &users)
		json.Unmarshal([]byte(data[i]["dates_json"]), &dates_json)

		for _, item := range users {

			if _, ok := result[item.User]; ok {

				for j, date := range dates_json {

					if _, ok := result[item.User]; ok {
						result[item.User][j] = map[string]string{}

						result[item.User][j]["date_start"] = date.DateStart
						result[item.User][j]["date_end"] = date.DateEnd

						if data[i]["is_external"] == "1" {

							result[item.User][j]["is_external"] = "Внешнее"

						} else {

							result[item.User][j]["is_external"] = "Внутреннее" //внутреннее

						}

						result[item.User][j]["topic_title"] = data[i]["topic_title"]
						result[item.User][j]["deps"] = usersindeps[item.User]

					} else {
						result[item.User][j]["date_start"] = date.DateStart
						result[item.User][j]["date_end"] = date.DateEnd

						if data[i]["is_external"] == "1" {

							result[item.User][j]["is_external"] = "Внешнее"

						} else {

							result[item.User][j]["is_external"] = "Внутреннее" //внутреннее

						}

						result[item.User][j]["topic_title"] = data[i]["topic_title"]
						result[item.User][j]["deps"] = usersindeps[item.User]

					}
				}

			} else {

				result[item.User] = map[int]map[string]string{}

				for j, date := range dates_json {

					if _, ok := result[item.User]; ok {
						result[item.User][j] = map[string]string{}

						result[item.User][j]["date_start"] = date.DateStart
						result[item.User][j]["date_end"] = date.DateEnd

						if data[i]["is_external"] == "1" {

							result[item.User][j]["is_external"] = "Внешнее"

						} else {

							result[item.User][j]["is_external"] = "Внутреннее" //внутреннее

						}

						result[item.User][j]["topic_title"] = data[i]["topic_title"]
						result[item.User][j]["deps"] = usersindeps[item.User]

					} else {
						result[item.User][j]["date_start"] = date.DateStart
						result[item.User][j]["date_end"] = date.DateEnd

						if data[i]["is_external"] == "1" {

							result[item.User][j]["is_external"] = "Внешнее"

						} else {

							result[item.User][j]["is_external"] = "Внутреннее" //внутреннее

						}

						result[item.User][j]["topic_title"] = data[i]["topic_title"]
						result[item.User][j]["deps"] = usersindeps[item.User]

					}
				}

			}

		}

	}
	return result

}

//GetExelAnaliticsTraining Get Exel Analitics Training
func GetExelAnaliticsTraining(c *gin.Context) {

	users := allUsersIntraining()
	_ = users

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   users,
	})
}
