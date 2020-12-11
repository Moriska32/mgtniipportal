package training

import (
	"PortalMGTNIIP/config"
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
		year_text = fmt.Sprintf(`EXTRACT(YEAR from (training.dates_json -> 0 ->> 'date_start')::date) = %s `, year)
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
	FROM public.training where  %s %s %s %s
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
