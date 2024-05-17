package get

import (
	. "masterservice/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func GetSessionHost(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})

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

	db.Conn.Debug().Model(&cur).Where("issession = true AND isactive = true").First(&cur)

	ans["host"] = cur.Host
	ans["port"] = cur.Port

	response = Interfacetoresponse(t, ans)
	return response
}
