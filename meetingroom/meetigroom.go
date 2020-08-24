package meetingroom

import (
	"PortalMGTNIIP/config"
	"fmt"
	"net/http"
	"strings"

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
