package post

import (
	"fmt"
	. "masterservice/model"
	"time"

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
func GetCronMandate(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["uuid"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing UUID")
	}

	uuid := p.Sanitize(fmt.Sprintf("%v", args["uuid"]))

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

	if args["masterid"] != nil {
		//Check does it is still master
		ms := &Microservices{}
		masterid := p.Sanitize(fmt.Sprintf("%v", args["masterid"]))
		rows := db.Conn.Where(`uuid = ? AND masterid = ? AND isactive = true`, uuid, masterid).First(&ms)
		if rows.RowsAffected == 0 {
			ans["mandate"] = false
			return Interfacetoresponse(t, ans)
		}

		//update curtime
		curtime := GetTime()
		newms := &Microservices{}
		newms.Echo = curtime
		err = db.Conn.Where(`uuid = ? AND masterid = ? AND isactive = true`, uuid, masterid).Updates(&newms).Error
		if err != nil {
			return ErrorReturn(t, 400, "000005", err.Error())
		}

		ans["mandate"] = true
		return Interfacetoresponse(t, ans)

	}

	//It is slave. Check Does master is alive
	ms := &Microservices{}
	rows := db.Conn.Where(`uuid = ? AND isactive = true`, uuid).First(&ms)
	if rows.RowsAffected == 0 {
		return ErrorReturn(t, 400, "000005", "Microservice not found or deactivated")
	}

	curtime := GetTime()
	lastechotime := ms.Echo
	df := curtime - lastechotime

	if df > 5 {
		//Grant master mandate
		newmasterid := Hashgen(24)

		newms := &Microservices{}
		newms.Echo = curtime
		newms.MasterID = newmasterid
		err = db.Conn.Where(`uuid = ? AND isactive = true`, uuid).Updates(&newms).Error
		if err != nil {
			return ErrorReturn(t, 400, "000005", err.Error())
		}

		ans["masterid"] = newmasterid
		ans["mandate"] = true
		return Interfacetoresponse(t, ans)
	}

	ans["mandate"] = false
	return Interfacetoresponse(t, ans)

}

func GetTime() int64 {

	now := time.Now()
	r := now.Unix()

	return r
}
