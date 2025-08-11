package post

import (
	"context"
	"fmt"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// перезагружает контейнер по имени или ID
func RestartContainer(containerName string) error {
	stopOptions := container.StopOptions{
		Timeout: nil,
	}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	cli.NegotiateAPIVersion(context.Background())

	ctx := context.Background()

	// Попытка перезапуска контейнера с таймаутом
	if err := cli.ContainerRestart(ctx, containerName, stopOptions); err != nil {
		return fmt.Errorf("failed to restart container %s: %w", containerName, err)
	}

	fmt.Printf("Container %s restarted successfully\n", containerName)
	return nil
}

func Restart(t *pb.Request) (response *pb.Response) {

	if *t.IsAdmin != 1 {
		response = ErrorReturn(t, 401, "000011", "You have no admin rights")
	}

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["container"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing  important data")
	}

	containerName := p.Sanitize(fmt.Sprintf("%v", args["container"]))

	if err := RestartContainer(containerName); err != nil {

		return ErrorReturn(t, 406, "000013", err.Error())
	}

	ans["answer"] = "done"
	response = Interfacetoresponse(t, ans)
	return response
}
