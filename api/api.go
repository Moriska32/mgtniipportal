package api

import (
	config "ProtalMGTNIIP/config"
	"log"
	"net/http"

	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
)

func Dep(c *gin.Context) {
	dbConnect := config.Connect()
	ID := c.Param("id")
	todo := "SELECT dep_id, name, parent_id FROM public.tdep where parent_id = " + ID + ";"
	rows, err := dbConnect.Query(todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data":    jsonify.Jsonify(rows),
	})
	return

}
