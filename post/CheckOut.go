package post

import (
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
func CheckOut(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["uuid"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing  UUID")
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

	newdata := &Microservices{}
	//newdata.IsActive = false

	err = db.Conn.Model(newdata).Where("uuid = ?", uuid).Update("isactive", false).Error
	if err != nil {
		return ErrorReturn(t, 400, "000005", err.Error())
	}

	ans["answer"] = "Done"
	response = Interfacetoresponse(t, ans)

	return response
}
