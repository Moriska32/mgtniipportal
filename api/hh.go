package api

import (
	"PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//HHJSON hh json
type HHJSON struct {
	DepID string `json:"dep_id"`
	Post  string `json:"post"`
	Descr string `json:"descr"`
	Vacid string `json:"vac_id"`
}

//PostHH add to BD
func PostHH(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	var json HHJSON
	c.BindJSON(&json)

	insert := `INSERT INTO public.tvacancy
	(dep_id, post, descr)
	VALUES(?, '?', '?');
	`
	_ = dbConnect.QueryRow(insert, json.DepID, json.Post, json.Descr)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
	return
}

//UpdateHH update hh in BD
func UpdateHH(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	var json HHJSON
	c.BindJSON(&json)

	insert := `UPDATE public.tvacancy
	SET dep_id=?, post='?', descr='?'
	WHERE vac_id=?;
	`
	_ = dbConnect.QueryRow(insert, json.DepID, json.Post, json.Descr, json.Vacid)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
	return
}

//DeleteHH Delete hh in BD
func DeleteHH(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	var json HHJSON
	c.BindJSON(&json)

	insert := `DELETE FROM public.tvacancy
	WHERE vac_id=?;
	`
	_ = dbConnect.QueryRow(insert, json.Vacid)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
	return
}

//GetHHs get hh in BD
func GetHHs(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := `SELECT vac_id, dep_id, post, descr
	FROM public.tvacancy;
	`

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

//GetHH get hh in BD
func GetHH(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()
	ID := c.Param("id")
	todo := fmt.Sprintf(`SELECT vac_id, dep_id, post, descr
	FROM public.tvacancy where vac_id = %s;
	`, ID)

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
	return
}
