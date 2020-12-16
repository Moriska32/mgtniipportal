package training

import (
	"PortalMGTNIIP/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
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

func allUsersIntraining() map[string]map[string]string {

	usersindeps := getdepbyuser()

	var result = map[string]map[string]string{}

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

				if _, ok := result[item.User]; ok {
					result[item.User] = map[string]string{}

					result[item.User]["date_start"] = dates_json[0].DateStart
					result[item.User]["date_end"] = dates_json[len(dates_json)-1].DateEnd

					if data[i]["is_external"] == "1" {

						result[item.User]["is_external"] = "Внешнее"

					} else {

						result[item.User]["is_external"] = "Внутреннее" //внутреннее

					}

					result[item.User]["topic_title"] = data[i]["topic_title"]
					result[item.User]["deps"] = usersindeps[item.User]

				} else {
					result[item.User]["date_start"] = dates_json[0].DateStart
					result[item.User]["date_end"] = dates_json[len(dates_json)-1].DateEnd

					if data[i]["is_external"] == "1" {

						result[item.User]["is_external"] = "Внешнее"

					} else {

						result[item.User]["is_external"] = "Внутреннее" //внутреннее

					}

					result[item.User]["topic_title"] = data[i]["topic_title"]
					result[item.User]["deps"] = usersindeps[item.User]

				}

			} else {

				result[item.User] = map[string]string{}

				if _, ok := result[item.User]; ok {
					result[item.User] = map[string]string{}

					result[item.User]["date_start"] = dates_json[0].DateStart
					result[item.User]["date_end"] = dates_json[len(dates_json)-1].DateEnd

					if data[i]["is_external"] == "1" {

						result[item.User]["is_external"] = "Внешнее"

					} else {

						result[item.User]["is_external"] = "Внутреннее" //внутреннее

					}

					result[item.User]["topic_title"] = data[i]["topic_title"]
					result[item.User]["deps"] = usersindeps[item.User]

				} else {
					result[item.User]["date_start"] = dates_json[0].DateStart
					result[item.User]["date_end"] = dates_json[len(dates_json)-1].DateEnd

					if data[i]["is_external"] == "1" {

						result[item.User]["is_external"] = "Внешнее"

					} else {

						result[item.User]["is_external"] = "Внутреннее" //внутреннее

					}

					result[item.User]["topic_title"] = data[i]["topic_title"]
					result[item.User]["deps"] = usersindeps[item.User]

				}

			}

		}

	}

	return result

}

func userwithname() map[string]map[string]string {

	result := map[string]map[string]string{}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`select tuser.user_id, tuser.fam, tuser."name", tuser.otch 
	FROM public.tuser;`)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)

	}

	for _, item := range data {

		if _, ok := result[item["user_id"]]; ok {

			result[item["user_id"]]["fam"] = item["fam"]
			result[item["user_id"]]["name"] = item["name"]
			result[item["user_id"]]["otch"] = item["otch"]

		} else {

			result[item["user_id"]] = map[string]string{}

			result[item["user_id"]]["fam"] = item["fam"]
			result[item["user_id"]]["name"] = item["name"]
			result[item["user_id"]]["otch"] = item["otch"]

		}

	}
	return result
}

func writeUsersrTrainingToExcel() {

	users := allUsersIntraining()
	name := userwithname()

	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "Фамилия")
	f.SetCellValue("Sheet1", "B1", "Имя")
	f.SetCellValue("Sheet1", "C1", "Отчество")
	f.SetCellValue("Sheet1", "D1", "Начало обучения")
	f.SetCellValue("Sheet1", "E1", "Конец обучения")
	f.SetCellValue("Sheet1", "F1", "Вид обучения")
	f.SetCellValue("Sheet1", "G1", "Тема обучения")
	f.SetCellValue("Sheet1", "H1", "Подразделение сотрудников")
	// Set active sheet of the workbook.

	i := 2

	for user, items := range users {

		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), name[user]["fam"])
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), name[user]["name"])
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), name[user]["otch"])
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), items["date_start"])
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i), items["date_end"])
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i), items["is_external"])
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i), items["topic_title"])
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i), items["deps"])

		i++

	}
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("public//documents//excel//Данные по пользователям по обучениям.xlsx"); err != nil {
		fmt.Println(err)
	}

}

//GetExelAnaliticsTraining Get Exel Analitics Training
func GetExelAnaliticsTraining(c *gin.Context) {

	writeUsersrTrainingToExcel()

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"file":   "172.20.0.82:4747/file/documents/excel/Данные по пользователям по обучениям.xlsx",
	})

}
