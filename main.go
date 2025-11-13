package main

import (
	"fmt"
	"masterservice/entrypoint"
	. "masterservice/global"
	. "masterservice/version"
	"net"
	"os"
	"strings"
	"time"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	. "github.com/gogufo/gufo-api-gateway/gufodao"

	"github.com/certifi/gocertifi"
	"github.com/getsentry/sentry-go"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

/*
var (

	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:5300", "gRPC server endpoint")

)
*/
func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("./config/")

	//port := viper.GetInt("server.port")
	//	portsrting := fmt.Sprintf(":%d", port)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error with Settings \n")
		os.Exit(3)
	}

	if viper.GetBool("server.sentry") {

		SetLog("Connect to Setry...")

		sentryClientOptions := sentry.ClientOptions{
			Dsn:              viper.GetString("sentry.dsn"),
			EnableTracing:    viper.GetBool("sentry.tracing"),
			Debug:            viper.GetBool("sentry.debug"),
			TracesSampleRate: viper.GetFloat64("sentry.trace"),
		}

		rootCAs, err := gocertifi.CACerts()
		if err != nil {
			SetLog("Could not load CA Certificates for Sentry: " + err.Error())

		} else {
			sentryClientOptions.CaCerts = rootCAs
		}

		err = sentry.Init(sentryClientOptions)

		if err != nil {
			SetLog("Error with sentry.Init: " + err.Error())
		}

		flushsec := viper.GetDuration("sentry.flush")

		defer sentry.Flush(flushsec * time.Second)

	}
	portvar := fmt.Sprintf("microservices.%s.port", MicroServiceName)
	getport := viper.GetString(portvar)
	port := ":5300"
	if getport != "" {
		port = fmt.Sprintf(":%s", getport)
	}

	listener, err := net.Listen("tcp", port)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s := &Server{}

	setingskey := fmt.Sprintf("microservices.%s.entrypointversion", MicroServiceName)
	lastentrypointversion := viper.GetString(setingskey)

	if lastentrypointversion != VERSIONPLUGIN {
		go entrypoint.Init()
	}

	pb.RegisterReverseServer(grpcServer, s)

	grpcServer.Serve(listener)

}

type Server struct {
}

func (s *Server) Do(c context.Context, request *pb.Request) (response *pb.Response, err error) {

	/// ==========================
	//  SECURITY CHECK (unified with Gufo Gateway)
	// ==========================
	mode := strings.ToLower(viper.GetString("security.mode"))

	switch mode {

	case "hmac":
		secret := viper.GetString("security.hmac_secret")
		maxAge := time.Duration(viper.GetInt("security.max_age")) * time.Second

		if request.Sign == nil || request.Module == nil ||
			!VerifyHMAC(secret, *request.Module, *request.Sign, maxAge) {

			return ErrorReturn(request, 401, "00001", "Invalid or expired HMAC signature"), nil
		}

	case "sign":
		if request.Sign == nil || viper.GetString("server.sign") != *request.Sign {
			return ErrorReturn(request, 401, "00001", "Invalid signature"), nil
		}

	case "mtls":
		// gRPC TLS already validated client certificate if configured
		SetLog("mTLS mode active - skipping sign verification")

	default:
		return ErrorReturn(request, 500, "00002", "Security mode not configured"), nil
	}

	//Check connection

	response = Init(request)

	return response, nil
}
