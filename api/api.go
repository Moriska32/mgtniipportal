package api

import (
	config "PortalMGTNIIP/config"
	"encoding/json"
	"io/ioutil"
	"log"

	"net/http"

	"github.com/elgs/gosqljson"

	"github.com/gin-gonic/gin"
)

//Dep list if dep by id
func Dep(c *gin.Context) {
	dbConnect := config.Connect()
	ID := c.Param("id")
	todo := "SELECT dep_id, name, parent_id FROM public.tdep where parent_id = " + ID + ";"

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

//Deps List all of deps
func Deps(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT dep_id, name, parent_id FROM public.tdep;"

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

//Posts List all of Post
func Posts(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT post_id, dep_id, post_name, post_count FROM public.tpost;"

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

//Post list if dep by id
func Post(c *gin.Context) {
	dbConnect := config.Connect()
	ID := c.Param("id")
	todo := "SELECT  post_id, dep_id, post_name, post_count FROM public.tdep where dep_id = " + ID + ";"

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

//Objectstype List all of Objectstype
func Objectstype(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT type_id, type_name, container FROM public.sobject_type;"

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

// Cbrdaily get values
func Cbrdaily(c *gin.Context) {

	url := "https://www.cbr-xml-daily.ru/daily_json.js"

	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	cbr, err := UnmarshalWelcome(body)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   cbr.Valute,
	})

}

//UnmarshalWelcome get values
func UnmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

//Marshal values
func (r *Welcome) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

//Welcome values
type Welcome struct {
	Date         string            `json:"Date"`
	PreviousDate string            `json:"PreviousDate"`
	PreviousURL  string            `json:"PreviousURL"`
	Timestamp    string            `json:"Timestamp"`
	Valute       map[string]Valute `json:"Valute"`
}

//Valute valute
type Valute struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  int64   `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}

//Weather get Weather
func Weather(c *gin.Context) {

	url := "https://accuweatherstefan-skliarovv1.p.rapidapi.com/get24HoursConditionsByLocationKey"

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("x-rapidapi-host", "AccuWeatherstefan-skliarovV1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "7e8e2a23afmshf2deabcb8c7b3dfp177c0cjsne0b1f5d1b74e")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   body,
	})

}
