package get

import (
	"fmt"
	. "masterservice/model"
	"strconv"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func List(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

	cur := []Microservices{}

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
		return ErrorReturn(t, 500, "000027", err.Error())
	}

	offset := 0
	limit := 25

	if args["offset"] != nil {
		offset, _ = strconv.Atoi(fmt.Sprintf("%v", args["offset"]))
	}

	if args["limit"] != nil {
		limit, _ = strconv.Atoi(fmt.Sprintf("%v", args["limit"]))
	}

	var count int64

	db.Conn.Debug().Model(&cur).Count(&count)
	db.Conn.Debug().Model(&cur).Limit(limit).Offset(offset).Find(&cur)

	ans["microservices"] = cur
	ans["mccount"] = count

	response = Interfacetoresponse(t, ans)
	return response
}
