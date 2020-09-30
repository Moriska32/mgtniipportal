package geomap

import (
	config "PortalMGTNIIP/config"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Geom List all of deps
type Geom struct {
	type_id      string
	type_name    string
	container    string
	object_id    string
	st_asgeojson string
	container_id sql.NullString
}

//Map get geom by floor
func Map(c *gin.Context) {

	geom := []*Geom{}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	ID := "1"

	sql := `SELECT sobject_type.*, tobject.object_id, St_asgeojson(tobject.geom), tobject.container_id
		FROM public.sobject_type sobject_type, public.tobject tobject
		WHERE
			tobject.type_id = sobject_type.type_id and tobject.container_id = %s;`

	_ = sql

	todo := `SELECT sobject_type.*, tobject.object_id, St_asgeojson(tobject.geom), tobject.container_id
		FROM public.sobject_type sobject_type, public.tobject tobject
		WHERE
			tobject.type_id = sobject_type.type_id and tobject.object_id = ` + ID + `;`

	fmt.Println(todo)

	rows, err := dbConnect.Query(todo)

	_ = rows

	for rows.Next() {
		pool := new(Geom)

		if err = rows.Scan(&pool.type_id, &pool.type_name, &pool.container, &pool.object_id, &pool.st_asgeojson, &pool.container_id); err != nil {
			fmt.Println("Scanning failed.....")
			fmt.Println(err.Error())
			return
		}

		geom = append(geom, pool)
		fmt.Print(geom)
		// rowson, err := dbConnect.Query(fmt.Sprintf(sql, pool.object_id))

		// for rowson.Next() {

		// 	if err = rowson.Scan(&pool.type_id, &pool.type_name, &pool.container, &pool.object_id, &pool.st_asgeojson, &pool.container_id); err != nil {
		// 		fmt.Println("Scanning rowson failed.....")
		// 		fmt.Println(err.Error())
		// 		continue
		// 	}
		// 	geom = append(geom, pool)

		// 	rowsons, err := dbConnect.Query(fmt.Sprintf(sql, pool.object_id))

		// 	for rowsons.Next() {

		// 		if err = rowsons.Scan(&pool.type_id, &pool.type_name, &pool.container, &pool.object_id, &pool.st_asgeojson, &pool.container_id); err != nil {
		// 			fmt.Println("Scanning rowsons failed.....")
		// 			fmt.Println(err.Error())
		// 			continue
		// 		}
		// 		geom = append(geom, pool)

		// 	}

		// }
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   geom,
	})

	return

}
