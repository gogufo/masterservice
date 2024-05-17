package main

import (
	gt "masterservice/get"
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

	switch *t.Method {
	case "GET":
		switch *t.Param {
		case "health":
			response = health(t)
		default:
			response = gt.Init(t)
		}
	case "POST":
		response = pt.Init(t)
	case "PATCH":
		response = pa.Init(t)

	}

	return response

}

func health(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["health"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}
