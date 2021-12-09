package carscannertodb

import (
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxClient() {
	bucket := "example-bucket"
	org := "example-org"
	token := "example-token"
	// Store the URL of your InfluxDB instance
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	api := client.WriteAPI(org, bucket)
	fmt.Println(api)
}
