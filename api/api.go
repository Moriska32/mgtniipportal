package api

import (
	config "PortalMGTNIIP/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

//Deps_id List all of deps
type Deps_id struct {
	Dep_id      int
	Name        string
	Parent_id   int
	Child_posts []*Posts_id
	Child_deps  []*Deps_id
}

//Posts_id List all of deps
type Posts_id struct {
	Post_id   int
	Dep_id    int
	Post_name string
}

func checkinstruct(pool []*Posts_id, post *Posts_id) bool {

	result := true

	for _, item := range pool {

		if item.Post_id == post.Post_id {
			result = false
			return result
		}

	}
	return result
}

func checkinstructdeps(pool []*Deps_id, post *Deps_id) bool {

	result := true

	for _, item := range pool {

		if item.Dep_id == post.Dep_id {
			result = false
			return result
		}

	}
	return result
}

//Orgstructure List all of deps
func Orgstructure(c *gin.Context) {

	deps := []*Deps_id{}

	dbConnect := config.Connect()
	todo := `SELECT dep_id, name, parent_id FROM public.tdep where dep_id != 1 and dep_id not in (3, 27, 29, 64, 67, 69);`

	rows, err := dbConnect.Query(todo)

	defer rows.Close()
	for rows.Next() {
		pool := new(Deps_id)

		if err = rows.Scan(&pool.Dep_id, &pool.Name, &pool.Parent_id); err != nil {
			fmt.Println("Scanning failed.....")
			fmt.Println(err.Error())
			return
		}

		deps = append(deps, pool)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	todo = `SELECT post_id, dep_id, post_name FROM public.tpost;`
	rows, err = dbConnect.Query(todo)
	posts := []*Posts_id{}
	defer rows.Close()
	for rows.Next() {
		pool := new(Posts_id)

		if err = rows.Scan(&pool.Post_id, &pool.Dep_id, &pool.Post_name); err != nil {
			fmt.Println("Scanning failed.....")
			fmt.Println(err.Error())
			return
		}

		posts = append(posts, pool)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	j := 0
	k := 0
	result := []*Deps_id{}
	for i, dep := range deps {
		if dep.Parent_id == 1 {
			result = append(result, dep)

			for _, post := range posts {

				if dep.Dep_id == post.Dep_id {
					result[i].Child_posts = append(result[i].Child_posts, post)

				}

			}

			for _, subdep := range deps {

				if result[i].Dep_id == subdep.Parent_id && checkinstructdeps(result[i].Child_deps, subdep) {

					result[i].Child_deps = append(result[i].Child_deps, subdep)

					for _, post := range posts {

						if result[i].Child_deps[j].Dep_id == post.Dep_id && checkinstruct(result[i].Child_deps[j].Child_posts, post) {

							result[i].Child_deps[j].Child_posts = append(result[i].Child_deps[j].Child_posts, post)

						}

					}

					for _, subsubdep := range deps {

						if result[i].Child_deps[j].Dep_id == subsubdep.Parent_id && checkinstructdeps(result[i].Child_deps[j].Child_deps, subsubdep) {

							result[i].Child_deps[j].Child_deps = append(result[i].Child_deps[j].Child_deps, subsubdep)

							for _, post := range posts {

								if result[i].Child_deps[j].Child_deps[k].Dep_id == post.Dep_id && checkinstruct(result[i].Child_deps[j].Child_deps[k].Child_posts, post) {

									result[i].Child_deps[j].Child_deps[k].Child_posts = append(result[i].Child_deps[j].Child_deps[k].Child_posts, post)

								}

							}
							k++
						}

					}
					k = 0
					j++

				}

			}
			j = 0
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   result,
	})
	dbConnect.Close()
	return

}

//Deps List all of deps
func Deps(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT dep_id, name, parent_id FROM public.tdep where dep_id not in (3, 27, 29, 64, 67, 69);"

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

	url := "http://api.openweathermap.org/data/2.5/weather?q=Moscow&appid=1e0cab77972f211e662fccf809bafc72"

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

//Meetingrooms get Meetingrooms
func Meetingrooms(c *gin.Context) {

	dbConnect := config.Connect()
	todo := `SELECT object_id, number
	FROM public.tobject where "number" in ('2','9','505') and type_id = 4;`

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

//Objects get Object
func Objects(c *gin.Context) {

	ID := c.Param("id")
	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := `SELECT object_id, type_id, container_id, "number"
	FROM public.tobject where type_id in (1,2,4) and number not in ('0', '') and  LENGTH(number) <= 5 and type_id = ` + ID + `;`

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
