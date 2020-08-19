package api

import (
	config "PortalMGTNIIP/config"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"net/http"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//Dep list if dep by id
func Dep(c *gin.Context) {
	dbConnect := config.Connect()
	ID := c.Param("id")
	todo := "SELECT dep_id, name, parent_id FROM public.tdep where parent_id = " + ID + " and parent_id not in (3, 27, 29, 64, 67, 69);"

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

//Orgstructurestruct List all of deps
type Orgstructurestruct []struct {
	Name      string `json:"name"`
	DepID     int    `json:"dep_id"`
	ChildDeps [][]struct {
		DepID      int    `json:"dep_id"`
		Name       string `json:"name"`
		ParentID   int    `json:"parent_id"`
		ChildPosts []struct {
			DepID     int    `json:"dep_id"`
			PostID    int    `json:"post_id"`
			PostName  string `json:"post_name"`
			PostCount int    `json:"post_count"`
		} `json:"child_posts"`
	} `json:"child_deps"`
}

//Orgstructure List all of deps
func Orgstructure(c *gin.Context) {

	dbConnect := config.Connect()
	todo := ` select replace(replace(replace(replace(
		(select jsonb_agg(result) as result from (
		select name, dep_id, child_deps from
		(select *,
		(select jsonb_agg(child_deps) from (select a.dep_id as dep_id, a.name as name, a.parent_id as parent_id, 
		 (select json_agg(child_deps) from
			(select * from 
			(select dep.dep_id, dep.name, dep.parent_id, 
			(select jsonb_agg(child_posts) from 
			(select * from public.tpost where dep_id = dep.dep_id) child_posts
			)::text as child_posts
		   from public.tdep dep) as posts
		  where child_posts is not null) child_deps
			 )::text as child_deps
		  from public.tdep as a) child_deps
		 where child_deps is not null and b.dep_id = child_deps.parent_id )::text as child_deps
		from public.tdep b) res
		where child_deps is not null) result)::text ,'\n',''),'\',''),'"[','['),']"',']') as data;`
	rows, err := dbConnect.Query(todo)

	var data string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&data)
		if err != nil {
			log.Fatal(err)
		}

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	res := Orgstructurestruct{}
	json.Unmarshal([]byte(strings.Replace(data, `\`, ``, 1)), &res)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   res,
	})
	dbConnect.Close()
	return

}

//Deps List all of deps
func Deps(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT dep_id, name, parent_id FROM public.tdep where parent_id not in (3, 27, 29, 64, 67, 69);"

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
	dbConnect.Close()
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
	dbConnect.Close()
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
	dbConnect.Close()
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

	t := time.Now()

	type Weatherget map[string]float64

	url := "https://gridforecast.com/api/v1/forecast/55.7631;37.6241/" + t.Format("200601021500") + "?api_token=fi83J3miGOyofI5D"

	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	var r Weatherget
	err = json.Unmarshal(body, &r)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   r,
	})

}

//Weathers get Weather
func Weathers(c *gin.Context) {

	url := "api.openweathermap.org/data/2.5/weather?lat=55.7631&lon=37.6241&appid=1e0cab77972f211e662fccf809bafc72"

	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   body,
	})

}
