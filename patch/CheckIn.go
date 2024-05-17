package patch

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
func CheckIn(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["uuid"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing UUID")
	}

	if args["port"] != nil || args["host"] != nil {

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

		if args["port"] != nil {
			newdata.Port = p.Sanitize(fmt.Sprintf("%v", args["port"]))
		}
		if args["host"] != nil {
			newdata.Host = p.Sanitize(fmt.Sprintf("%v", args["host"]))
		}

		err = db.Conn.Model(&newdata).Where("uuid = ?", uuid).Updates(&newdata).Error
		if err != nil {
			return ErrorReturn(t, 400, "000005", err.Error())
		}

		ans["uuid"] = uuid
		response = Interfacetoresponse(t, ans)

		return response

	}

	return ErrorReturn(t, 400, "000007", "Missing Port or Host")

}
