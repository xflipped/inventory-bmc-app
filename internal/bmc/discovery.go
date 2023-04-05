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
	"github.com/foliagecp/inventory-bmc-app/sdk/pbdiscovery"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"github.com/stmcginnis/gofish"
	"go.mongodb.org/mongo-driver/bson"
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
