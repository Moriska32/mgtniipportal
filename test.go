package main

import (
	config "ProtalMGTNIIP/config"
	"fmt"
	_ "log"
	_ "os"

	"github.com/bdwilliams/go-jsonify/jsonify"
)

func maintest() {
	var print = fmt.Println
	dbConnect := config.Connect()

	ID := "1"

	todo := "SELECT dep_id, name, parent_id FROM public.tdep where parent_id = " + ID + ";"

	rows, err := dbConnect.Query(todo)

	if err != nil {
		panic(err.Error())
	}

	print(jsonify.Jsonify(rows))
}
