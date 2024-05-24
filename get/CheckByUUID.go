package get

import (
	. "masterservice/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func CheckByUUID(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})

	uuid := *t.ParamID

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
		return ErrorReturn(t, 500, "000027", err.Error())
	}

	ans["answer"] = true

	mcur := &Microservices{}
	rows := db.Conn.Debug().Model(&mcur).Where("uuid = ? AND isactive = true", uuid).First(&mcur)
	if rows.RowsAffected == 0 {
		ans["answer"] = false
		ans["name"] = mcur.Name
		ans["port"] = mcur.Port
		ans["host"] = mcur.Host
		ans["isinternal"] = mcur.IsInternal

	}

	response = Interfacetoresponse(t, ans)
	return response
}
