package geomap

import (
	config "PortalMGTNIIP/config"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Geom List all of deps
type Geom struct {
	TypeID      string
	TypeName    string
	Container   string
	ObjectID    string
	Geom        string
	ContainerID sql.NullString
}

//Map get geom by floor
func Map(c *gin.Context) {

	id := c.PostForm("id")

	geom := []*Geom{}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)

	fmt.Printf(string(bodyBytes))

	sql := `SELECT sobject_type.*, tobject.object_id, St_asgeojson(tobject.geom), tobject.container_id
		FROM public.sobject_type sobject_type, public.tobject tobject
		WHERE
			tobject.type_id = sobject_type.type_id and tobject.container_id = %s;`

	_ = sql

	todo := `SELECT sobject_type.*, tobject.object_id, St_asgeojson(tobject.geom), tobject.container_id
		FROM public.sobject_type sobject_type, public.tobject tobject
		WHERE
			tobject.type_id = sobject_type.type_id and tobject.object_id = ` + id + `;`

	fmt.Println(todo)

	rows, err := dbConnect.Query(todo)

	_ = rows

	for rows.Next() {
		pool := new(Geom)

		if err = rows.Scan(&pool.TypeID, &pool.TypeName, &pool.Container, &pool.ObjectID, &pool.Geom, &pool.ContainerID); err != nil {
			fmt.Println("Scanning failed.....")
			fmt.Println(err.Error())
			return
		}

		geom = append(geom, pool)

		rowson, err := dbConnect.Query(fmt.Sprintf(sql, pool.ObjectID))

		for rowson.Next() {

			if err = rowson.Scan(&pool.TypeID, &pool.TypeName, &pool.Container, &pool.ObjectID, &pool.Geom, &pool.ContainerID); err != nil {
				fmt.Println("Scanning rowson failed.....")
				fmt.Println(err.Error())
				continue
			}
			geom = append(geom, pool)

			rowsons, err := dbConnect.Query(fmt.Sprintf(sql, pool.ObjectID))

			for rowsons.Next() {

				if err = rowsons.Scan(&pool.TypeID, &pool.TypeName, &pool.Container, &pool.ObjectID, &pool.Geom, &pool.ContainerID); err != nil {
					fmt.Println("Scanning rowsons failed.....")
					fmt.Println(err.Error())
					continue
				}
				geom = append(geom, pool)

			}

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   geom,
	})

	return

}
