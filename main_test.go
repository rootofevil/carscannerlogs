package carscannertodb

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func Test_lineToData(t *testing.T) {
	type args struct {
		line      string
		delimiter string
		date      time.Time
	}
	tests := []struct {
		name string
		args args
		want CarData
	}{
		// TODO: Add test cases.
	}
	date, _ := time.Parse("2006-01-02", "2021-11-19")
	fullDate, _ := time.Parse("2006-01-02 15:04:05.999", "2021-11-19 16:45:38.072")
	tests = append(tests,
		struct {
			name string
			args args
			want CarData
		}{
			name: "simple",
			args: args{line: "\"60338.071981\";\"Обороты двигателя\";\"994\";\"rpm\"", delimiter: ";", date: date},
			want: CarData{Second: 60338.071981,
				Pid:   "Обороты двигателя",
				Value: 994,
				Units: "rpm",
				Time:  fullDate},
		},
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lineToData(tt.args.line, tt.args.delimiter, tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lineToData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadCsv(t *testing.T) {
	readCsv("/Users/roe/CarLogs/CarScanner/2021-11-19 16-45-25.csv", ";")
}

func TestMain(m *testing.M) {
	token := "F-QFQpmCL9UkR3qyoXnLkzWj03s6m4eCvYgDl1ePfHBf9ph7yxaSgQ6WN0i9giNgRTfONwVMK1f977r_g71oNQ=="
	// Store the URL of your InfluxDB instance
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	fmt.Println(client.Options().SetPrecision(time.Millisecond))
	defer client.Close()
	dataset, _ := readCsv("/Users/roe/CarLogs/CarScanner/2021-11-19 16-45-25.csv", ";")
	dataset = dataset[0:2]
	for _, d := range dataset {
		fmt.Printf("%+v\n", d)
		d.SendToInfluxDb(client, "first", "audi")
		r, err := client.QueryAPI("first").Query(context.Background(),
			fmt.Sprintf(
				`from(bucket:"audi")
		|> range(start: -30d)
		|> filter(fn: (r) => r._measurement == "%s")`,
				d.Pid))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", r)
	}
}
