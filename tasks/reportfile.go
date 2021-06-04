package tasks

import (
	"PortalMGTNIIP/config"
	"fmt"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//BuildReport get pool from bd in to file
func BuildReport(c *gin.Context) {

	execute_start_plan_time := c.Query("execute_start_plan_time")
	execute_end_plan_time := c.Query("execute_end_plan_time")

	execute_start_time := c.Query("execute_start_time")
	execute_end_time := c.Query("execute_end_time")
	mark := c.Query("mark")
	status := c.Query("status")

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`SELECT id, type_id, "number", description, author_id, 
	(select fam from tuser where user_id = author_id) as fam,
	(select "name" from tuser where user_id = author_id) as "name",
	(select "otch" from tuser where user_id = author_id) as "otch",
	(select case when execute_end_plan_time > execute_end_time then 'X'
	when execute_end_plan_time < execute_end_time then 'V'
	else ''
	end as mark),
	(select case when execute_end_time is not null and (execute_decline_time is null and operator_decline_time is null) then 'В работе'
	when execute_end_time is null and  (execute_decline_time is null and operator_decline_time is null) then 'Выполнена'
	when  execute_decline_time is not null or operator_decline_time is not null then 'Отклонена'
	end as status),
	(select fam from tuser where user_id = executor_id) as exfam,
	(select "name" from tuser where user_id = executor_id) as exname,
	(select "otch" from tuser where user_id = executor_id) as exotch,
	phone, operator_id, executor_id, to_char(create_time,'YYYY-MM-DD HH:MI:SS') as create_time, operator_accept_time, operator_decline_time,
	to_char(execute_start_time,'YYYY-MM-DD HH:MI:SS') as execute_start_time, 
	to_char(execute_end_time,'YYYY-MM-DD HH:MI:SS') as execute_end_time,
	to_char(execute_start_plan_time,'YYYY-MM-DD HH:MI:SS') as execute_start_plan_time,
	to_char(execute_end_plan_time,'YYYY-MM-DD HH:MI:SS') as execute_end_plan_time,
	operator_comment,
	executor_comment, 
	to_char(execute_accept_time,'YYYY-MM-DD HH:MI:SS') as execute_accept_time,
	to_char(execute_decline_time,'YYYY-MM-DD HH:MI:SS') as execute_decline_time
	FROM public.tasks where id != '0'`)

	switch {

	case execute_start_plan_time != "" && execute_end_plan_time != "":
		sql = sql + "and execute_start_plan_time >= " + execute_start_plan_time + " and execute_end_plan_time <= " + execute_end_plan_time
	case mark != "":
		sql = sql + " and mark = " + mark
	case status != "":
		sql = sql + " and status = " + status
	case execute_start_time != "" && execute_end_time != "":
		sql = sql + "and execute_start_time >= " + execute_start_time + " and execute_end_time <= " + execute_end_time
	}

	sql = sql + `;`

	theCase := "lower"
	pool, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	if err != nil {
		fmt.Println(err)
	}
	f, err := excelize.OpenFile("tasks/Отчёт портал.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	currentTime := time.Now()
	for items := range pool {

		f.SetCellValue("report", fmt.Sprintf("A%d", 3+items), pool[items]["number"])
		f.SetCellValue("report", fmt.Sprintf("B%d", 3+items), fmt.Sprintf("%s %s %s", pool[items]["fam"], pool[items]["name"], pool[items]["otch"]))
		f.SetCellValue("report", fmt.Sprintf("C%d", 3+items), pool[items]["description"])
		f.SetCellValue("report", fmt.Sprintf("D%d", 3+items), pool[items]["execute_start_plan_time"])
		f.SetCellValue("report", fmt.Sprintf("E%d", 3+items), pool[items]["execute_end_plan_time"])
		f.SetCellValue("report", fmt.Sprintf("F%d", 3+items), pool[items]["execute_start_time"])
		f.SetCellValue("report", fmt.Sprintf("G%d", 3+items), pool[items]["execute_end_time"])
		f.SetCellValue("report", fmt.Sprintf("H%d", 3+items), pool[items]["mark"])
		f.SetCellValue("report", fmt.Sprintf("I%d", 3+items), pool[items]["status"])
		f.SetCellValue("report", fmt.Sprintf("I%d", 3+items), pool[items]["status"])
		f.SetCellValue("report", fmt.Sprintf("J%d", 3+items), fmt.Sprintf("%s %s %s", pool[items]["exfam"], pool[items]["exname"], pool[items]["exotch"]))
		f.SetCellValue("report", fmt.Sprintf("K%d", 3+items), pool[items]["executor_comment"])
	}
	filename := fmt.Sprintf("Отчёт_портал_%s.xlsx", currentTime.Format("2006.01.02 15_04_05"))
	if err := f.SaveAs(fmt.Sprintf("tasks/reports/%s", filename)); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"file":   "http://172.20.0.82:4747/tasksreport/" + filename,
		"data":   pool,
	})

	return

}
