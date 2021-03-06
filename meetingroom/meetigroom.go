package meetingroom

import (
	"PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//Newmeet Reserve meeting room
func Newmeet(c *gin.Context) {

	dbConnect := config.Connect()

	datetimes := c.PostFormArray("datetimes")
	objectid := c.PostForm("object_id")
	userid := c.PostForm("user_id")
	descr := c.PostForm("descr")

	for _, datetime := range datetimes {

		datebegin := strings.Split(datetime, "|")[0]
		dateend := strings.Split(datetime, "|")[1]

		todo := fmt.Sprintf(`INSERT INTO public.tobject_reserve
	(object_id, period_beg, period_end, user_id, descr)
	VALUES(%s, '%s', '%s', %s, '%s');`, objectid, datebegin, dateend, userid, descr)

		_, err := dbConnect.Exec(todo)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
		}

	}
	dbConnect.Close()
}

//Getmeets get all meeting by month
func Getmeets(c *gin.Context) {

	month := c.PostForm("month")
	year := c.PostForm("year")
	dbConnect := config.Connect()

	todo := fmt.Sprintf(`SELECT period_id, object_id, to_char(period_beg, 'YYYY-MM-DD HH24:MI') as period_beg, 
	to_char(period_end, 'YYYY-MM-DD HH24:MI') as period_end, user_id, descr
	FROM public.tobject_reserve 
	where extract(month from  period_beg) = %s and extract(year from  period_beg) = %s order by period_beg desc;`, month, year)

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

	dbConnect.Close()
}

//Deletemeet get all meeting by month
func Deletemeet(c *gin.Context) {

	periodid := c.PostForm("period_id")
	dbConnect := config.Connect()

	todo := fmt.Sprintf(`DELETE FROM public.tobject_reserve WHERE period_id = %s;`, periodid)

	theCase := "lower"
	_, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

	dbConnect.Close()
}

//Updatemeet all meeting by month
func Updatemeet(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()
	datetime := c.PostForm("datetime")
	objectid := c.PostForm("object_id")
	userid := c.PostForm("user_id")
	descr := c.PostForm("descr")
	periodid := c.PostForm("period_id")

	datebegin := strings.Split(datetime, "|")[0]
	dateend := strings.Split(datetime, "|")[1]

	todo := fmt.Sprintf(`UPDATE public.tobject_reserve
		SET object_id=%s, period_beg='%s', period_end='%s', user_id=%s, descr='%s'
		WHERE period_id= %s;`, objectid, datebegin, dateend, userid, descr, periodid)

	_, err := dbConnect.Exec(todo)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

}

//GetAllMeets get all meeting
func GetAllMeets(c *gin.Context) {

	dbConnect := config.Connect()

	todo := fmt.Sprintf(`SELECT period_id, object_id, to_char(period_beg, 'YYYY-MM-DD HH24:MI') as period_beg, 
	to_char(period_end, 'YYYY-MM-DD HH24:MI') as period_end, user_id, descr
	FROM public.tobject_reserve order by period_beg desc;`)

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

	dbConnect.Close()
}

//GetMeetsLimit get all meeting
func GetMeetsLimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	dbConnect := config.Connect()

	todo := fmt.Sprintf(`SELECT period_id, object_id, to_char(period_beg, 'YYYY-MM-DD HH24:MI') as period_beg, 
	to_char(period_end, 'YYYY-MM-DD HH24:MI') as period_end, user_id, descr
	FROM public.tobject_reserve order by period_beg desc limit %s offset %s;`, limit, offset)

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
	FROM public.tobject_reserve ;`, limit)

	count, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"count":  count[0]["pages_length"],
	})

	dbConnect.Close()
}
