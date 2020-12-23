package absence

import (
	"PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//PostAbsence Post Absence
func PostAbsence(c *gin.Context) {

	data := jwt.ExtractClaims(c)

	if data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	user_id := c.PostForm("user_id")
	date_start := c.PostForm("date_start")
	date_end := c.PostForm("date_end")
	absence_reason_id := c.PostForm("absence_reason_id")

	sql := fmt.Sprintf(`INSERT INTO public.absence
	(user_id, date_start, date_end, absence_reason_id)
	VALUES(%s, '%s', '%s', %s);
	`, user_id, date_start, date_end, absence_reason_id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//GetAbsencesLimit Get Absences Limit
func GetAbsencesLimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT *
	FROM public.absence	
	limit %s offset %s ;
	`, limit, offset)

	log.Println(todo)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	todo = fmt.Sprintf(`SELECT ceil(count(*)::real/%s::real) as pages_length
	FROM public.absence;
	`, limit)

	count, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"data":         data,
		"pages_length": count[0]["pages_length"],
	})

	return

}

//GetAbsencesMonth Get Absences Month
func GetAbsencesMonth(c *gin.Context) {

	year := c.Query("year")
	month := c.Query("month")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT absence_id, user_id, date_start, date_end, absence_reason_id
	FROM public.absence where 
	'%s-%s-01'::date between CONCAT(EXTRACT(YEAR from date_start::date)::text,'-',EXTRACT(month from date_start::date)::text,'-01')::date
	and 
	CONCAT(EXTRACT(YEAR from date_end::date)::text,'-',EXTRACT(month from date_end::date)::text,'-01')::date;
	`, year, month)

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

//GetAbsenceReasons Get Absence Reasons
func GetAbsenceReasons(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT absence_reason_id, absence_reason, color, bg
	FROM public.absence_reason;
	`)

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

//DeleteAbsence Delete Absence
func DeleteAbsence(c *gin.Context) {

	data := jwt.ExtractClaims(c)

	if data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	absence_id := c.PostForm("absence_id")

	sql := fmt.Sprintf(`DELETE FROM public.absence
	WHERE absence_id = %s
	
	`, absence_id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//UpdateAbsence Update Absence
func UpdateAbsence(c *gin.Context) {

	data := jwt.ExtractClaims(c)

	if data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	}

	user_id := c.PostForm("user_id")
	date_start := c.PostForm("date_start")
	date_end := c.PostForm("date_end")
	absence_reason_id := c.PostForm("absence_reason_id")
	absence_id := c.PostForm("absence_id")

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`UPDATE public.absence
	SET user_id=%s, date_start='%s', date_end='%s', absence_reason_id=%s
	where absence_id = %s;
	
	`, user_id, date_start, date_end, absence_reason_id, absence_id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}
