// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbdiscovery"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"github.com/stmcginnis/gofish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// device - remove if re-discovery and uuid updated
func (b *BmcApp) Discovery(ctx context.Context, request *pbdiscovery.Request) (device *pbredfish.Device, err error) {
	const colName = "devices"

	log.Infof("exec discovery: %s", request.GetUrl())

	u, err := url.Parse(request.GetUrl())
	if err != nil {
		return
	}

	host := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	defer httpClient.CloseIdleConnections()

	resp, err := httpClient.Get(fmt.Sprintf("%s/redfish/v1", host))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var gofishService gofish.Service
	if err = json.NewDecoder(resp.Body).Decode(&gofishService); err != nil {
		return
	}

	redfishDevice := db.RedfishDevice{
		Url:  host,
		UUID: gofishService.UUID,
	}

	filter := bson.D{{Key: "UUID", Value: redfishDevice.UUID}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishDevice); err != nil {
		return
	}

	return redfishDevice.ToProto()
}

func (b *BmcApp) ListDevices(ctx context.Context, empty *pbbmc.Empty) (devices *pbredfish.Devices, err error) {
	const colName = "devices"

	devicesCollection := b.database.Collection(colName)

	cur, err := devicesCollection.Find(ctx, bson.D{}, options.Find())
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	devices = &pbredfish.Devices{}

	var results []db.RedfishDevice

	if err = cur.All(ctx, &results); err != nil {
		return
	}

	for _, result := range results {
		device, err := result.ToProto()
		if err != nil {
			return devices, err
		}

		devices.Items = append(devices.Items, device)
	}

	return
}
