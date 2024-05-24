package main

import (
	gt "masterservice/get"
	. "masterservice/global"
	pa "masterservice/patch"
	pt "masterservice/post"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

/*
Endpoints
- Registration/Deregistration microservices
- Return microservice location (host and port)
- Cron balancer:
   - Master node should send echo
	 - Slave node sould check echo
	 - If echo is missing, masterservice grant master permissions

	 Registrtion information
	 - host
	 - port
	 - microservice name
	 - microservice group
	 - is it session holder?
	 - is itinteranal use Only

*/

func Init(t *pb.Request) (response *pb.Response) {

	method := *t.Method
	param := *t.Param

	if *t.Module != MicroServiceName {
		method = *t.IR.Method
		param = *t.IR.Param
	}

	switch method {
	case "GET":
		switch param {
		case "health":
			response = health(t)
		default:
			response = gt.Init(t)
		}
	case "POST":
		response = pt.Init(t)
	case "PATCH":
		response = pa.Init(t)
	default:
		response = ErrorReturn(t, 404, "000012", "Missing argument")

	}

	return response

}

func health(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["health"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}
