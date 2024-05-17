package get

import (
	. "masterservice/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func GetMSByPath(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})

	paramid := *t.ParamID
	param := *t.Param

	cur := &Microservices{}

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
		return ErrorReturn(t, 500, "000027", err.Error())
	}

	if t.ParamID != nil {
		mcur := &Microservices{}
		rows := db.Conn.Debug().Model(&mcur).Where("name = ? AND isactive = true AND group = ?", paramid, param).First(&mcur)
		if rows.RowsAffected == 0 {
			db.Conn.Debug().Model(&cur).Where("name = ? AND isactive = true", param).First(&cur)
		} else {
			cur.Host = mcur.Host
			cur.Port = mcur.Port
		}
	} else {
		db.Conn.Debug().Model(&cur).Where("name = ? AND isactive = true", param).First(&cur)
	}

	ans["host"] = cur.Host
	ans["port"] = cur.Port

	response = Interfacetoresponse(t, ans)
	return response
}
