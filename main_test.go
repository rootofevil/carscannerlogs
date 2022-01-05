package carscannertodb

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var testdatapath string = "test\\2021-11-19 16-45-25.csv"

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
	type args struct {
		path  string
		delim string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{name: "Windows Style Path", args: args{path: "test\\2021-11-19 16-45-25.csv", delim: ";"}, want: nil},
		{name: "Linux Style Path", args: args{path: "test/2021-11-19 16-45-25.csv", delim: ";"}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ReadCsv(tt.args.path, tt.args.delim); err != nil {
				t.Errorf("error = %v, want %v", err, tt.want)
			} else {
				t.Logf("Test %s, no errors", tt.name)
			}
		})
	}
	// _, err := ReadCsv(testdatapath, ";")
	// if err != nil {
	// 	t.Fail()
	// }
	// testdatapath = "test/2021-11-19 16-45-25.csv"
	// _, err := ReadCsv(testdatapath, ";")
	// if err != nil {
	// 	t.Fail()
	// }
}

func TestMain(t *testing.T) {
	token := "F-QFQpmCL9UkR3qyoXnLkzWj03s6m4eCvYgDl1ePfHBf9ph7yxaSgQ6WN0i9giNgRTfONwVMK1f977r_g71oNQ=="
	// Store the URL of your InfluxDB instance
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	fmt.Println(client.Options().SetPrecision(time.Millisecond))
	defer client.Close()
	dataset, err := ReadCsv(testdatapath, ";")
	if err != nil {
		t.Fail()
	}
	// dataset = dataset[0:2]
	for _, d := range dataset {
		fmt.Printf("%+v\n", d)
		d.SendToInfluxDb(client, "first", "audi")
		// 	query := fmt.Sprintf(
		// 		`from(bucket: "audi")
		// |> range(start: -30d)
		// |> filter(fn: (r) => r._measurement == "%s")`,
		// 		d.Pid)
		// 	r, err := client.QueryAPI("first").Query(context.Background(), query)
		// 	if err == nil {
		// 		// Iterate over query response
		// 		for r.Next() {
		// 			// Notice when group key has changed
		// 			if r.TableChanged() {
		// 				fmt.Printf("table: %s\n", r.TableMetadata().String())
		// 			}
		// 			// Access data
		// 			fmt.Printf("value: %v\n", r.Record().Value())
		// 		}
		// 		// Check for an error
		// 		if r.Err() != nil {
		// 			fmt.Printf("query parsing error: %s\n", r.Err().Error())
		// 		}
		// 	} else {
		// 		panic(err)
		// 	}
	}
}
