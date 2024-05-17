package post

import (
	"encoding/json"
	"fmt"
	. "masterservice/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

/*
Registrtion information
- host
- port
- microservice name
  - microservice group

- is it session holder?
- is itinteranal use Only
*/
func CheckIn(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["name"] == nil || args["host"] == nil || args["port"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing  important data")
	}

	name := p.Sanitize(fmt.Sprintf("%v", args["name"]))

	//1. Check does such miscoservice is exist

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	dbtype := viper.GetString("database.type")

	if !db.Conn.Migrator().HasTable(&Microservices{}) {
		//2. Create new table
		if dbtype == "mysql" {
			db.Conn.Debug().Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&Microservices{})
		} else {
			db.Conn.Debug().Migrator().CreateTable(&Microservices{})
		}
	}

	if args["group"] != nil {
		//Check does such group is exist
		nm := p.Sanitize(fmt.Sprintf("%v", args["group"]))
		ms := &Microservices{}
		rows := db.Conn.Where(`name = ?`, nm).First(&ms)
		if rows.RowsAffected == 0 {
			return ErrorReturn(t, 500, "000013", "No group found")
		}

		rows = db.Conn.Where(`name = ? AND isactive = true AND group = ?`, name, nm).First(&ms)
		if rows.RowsAffected != 0 {
			return ErrorReturn(t, 500, "000013", "Microservice with such id is exist. If you want to use this micreservice change the name or check out previous microservice")
		}

	} else {
		microserivces := &Microservices{}

		rows := db.Conn.Where(`name = ? AND isactive = true`, name).First(&microserivces)
		if rows.RowsAffected != 0 {
			return ErrorReturn(t, 500, "000013", "Microservice with such id is exist. If you want to use this micreservice change the name or check out previous microservice")
		}
	}

	if args["issession"] != nil {
		microserivces := &Microservices{}
		rows := db.Conn.Where(`name = ? AND isactive = true AND issession = true`, name).First(&microserivces)
		if rows.RowsAffected != 0 {
			return ErrorReturn(t, 500, "000013", "Microservice with such id is exist. If you want to use this micreservice change the name or check out previous microservice")
		}
	}

	newdata := &Microservices{}

	JsonArgs, err := json.Marshal(args)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}
	//2 Put args to Struct
	err = json.Unmarshal(JsonArgs, &newdata)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	newdata.UUID = Hashgen(24)
	newdata.IsActive = true

	err = db.Conn.Create(&newdata).Error
	if err != nil {
		return ErrorReturn(t, 400, "000005", err.Error())
	}

	ans["uuid"] = newdata.UUID
	response = Interfacetoresponse(t, ans)

	return response
}
