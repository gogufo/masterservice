package get

import (
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

/List
/gethostbyname
/getsessionhost
/
*/

func Init(t *pb.Request) (response *pb.Response) {

	switch *t.Param {
	case "list":
		response = List(t)
	case "getmicroservicebypath":
		response = GetMSByPath(t)
	case "getsessionhost":
		response = GetSessionHost(t)
	default:
		response = ErrorReturn(t, 404, "000012", "Missing argument")
	}

	return response

}
