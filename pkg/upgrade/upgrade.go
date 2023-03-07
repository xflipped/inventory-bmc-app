// Copyright 2023 NJWS Inc.

package upgrade

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
)

var (
	log = logrus.New()
)

type UpdateParameters struct {
	Targets []string
}

type Oem struct {
	ImageType string
}

func Upgrade(ctx context.Context, query, file, ftype, target string) (err error) {
	nodes, err := qdsl.Qdsl(ctx, query, qdsl.WithId(), qdsl.WithType())
	if err != nil {
		return
	}

	var wg sync.WaitGroup

	for _, node := range nodes {
		node := node
		log.Infof("document: %s", node.Id)

		if node.Type != types.RedfishDeviceKey {
			log.Infof("document: %s, skip...", node.Id)
			continue
		}

		var device device.RedfishDevice
		if err = json.Unmarshal(node.Object, &device); err != nil {
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := upgrade(ctx, device, file, ftype, target); err != nil {
				log.Error(err)
			}
		}()

	}

	wg.Wait()

	return
}

func upgrade(ctx context.Context, device device.RedfishDevice, file, ftype, target string) (err error) {
	u, err := url.Parse(device.Api)
	if err != nil {
		return
	}

	config := gofish.ClientConfig{
		Endpoint:  fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		Username:  device.Login,
		Password:  device.Password,
		Insecure:  true,
		BasicAuth: true,
	}

	client, err := gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}
	defer client.Logout()

	service := client.GetService()

	updateService, err := service.UpdateService()
	if err != nil {
		return
	}

	updateParameters := &UpdateParameters{
		Targets: []string{target},
	}

	updateParametersData, err := json.Marshal(updateParameters)
	if err != nil {
		return
	}

	oemData, err := json.Marshal(Oem{ImageType: ftype})
	if err != nil {
		return
	}

	file, err = filepath.Abs(file)
	if err != nil {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		return
	}

	payloadMap := map[string]io.Reader{
		"UpdateParameters": bytes.NewReader(updateParametersData),
		"OemParameters":    bytes.NewReader(oemData),
		"UpdateFile":       f,

		"@Redfish.OperationApplyTime": bytes.NewReader([]byte("Immediate")),
	}

	client, err = gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}

	resp, err := client.PostMultipart(updateService.MultipartHTTPPushURI, payloadMap)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println(device.Api, string(data))

	// taskService, err := service.TaskService()
	// if err != nil {
	// 	return
	// }

	// for {
	// 	tasks, err := taskService.Tasks()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if len(tasks) > 0 {
	// 		task := tasks[len(tasks)-1]

	// 		fmt.Printf("task %s (%d%%) state: %s (%s -> %s)\n", task.ID, task.PercentComplete, task.TaskState, task.StartTime, task.EndTime)

	// 	}

	// 	time.Sleep(time.Second)
	// }

	return
}
