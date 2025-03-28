package main

import (
	"context"
	"fmt"
	"github.com/project-kessel/inventory-client-go/common"
	"google.golang.org/protobuf/types/known/structpb"

	kesselv2 "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta2"
	v1beta2 "github.com/project-kessel/inventory-client-go/v1beta2"
)

func main() {

	data, _ := structpb.NewStruct(map[string]interface{}{
		"workspace_id": "ExampleResource",
	})

	data2, _ := structpb.NewStruct(map[string]interface{}{
		"satellite_id":          "21b4d79e-34ec-441d-8018-5f8985bb0413",
		"sub_manager_id":        "21b4d79e-34ec-441d-8018-5f8985bb0413",
		"insights_inventory_id": "21b4d79e-34ec-441d-8018-5f8985bb0413",
		"ansible_host":          "abddc",
	})

	client, err := v1beta2.NewHttpClient(context.Background(),
		common.NewConfig(common.WithHTTPUrl("localhost:8000")))
	common.WithTLSInsecure(true)
	// v1beta1.WithAuthEnabled("svc-test", "", "http://localhost:8084/realms/redhat-external/protocol/openid-connect/token"),
	// v1beta1.WithHTTPTLSConfig(tls.Config{})
	if err != nil {
		fmt.Println(err)
	}
	request := kesselv2.ReportResourceRequest{
		Resource: &kesselv2.Resource{
			ResourceType: "host",
			ReporterData: &kesselv2.ReporterData{
				ReporterType:       "HBI",
				ReporterInstanceId: "1",
				ReporterVersion:    "0.1",
				LocalResourceId:    "1",
				ApiHref:            "www.example.com",
				ConsoleHref:        "www.example.com",
				ResourceData:       data2,
			},
			CommonResourceData: data,
		},
	}
	// optts, err := client.GetTokenHTTPOption()
	resp, err := client.KesselResourceService.ReportResource(context.Background(), &request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
