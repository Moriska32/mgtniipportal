package api

import (
	"PortalMGTNIIP/config"
	"log"
	"net/http"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//Search in bd
func Search(c *gin.Context) {
	log.Println("поиск")

	param := c.PostForm("param")

	dbConnect := config.Connect()

	defer dbConnect.Close()

	log.Println(param)
	if param == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	refresh := `REFRESH MATERIALIZED VIEW public.searchlist;`

	_, err := dbConnect.Query(refresh)

	todo := `select pool from (
	SELECT row_to_json::text as pool
	FROM public.searchlist) as pool where pool like '%` + param + `%';`

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