package carscannertodb

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type CarData struct {
	Time   time.Time
	Second float64
	Pid    string
	Value  float64
	Units  string
}

func (cd CarData) SendToInfluxDb(client influxdb2.Client, org, bucket string) {
	api := client.WriteAPI(org, bucket)
	p := influxdb2.NewPoint(cd.Pid, map[string]string{"unit": cd.Units}, map[string]interface{}{"value": cd.Value}, time.Now())
	api.WritePoint(p)
}

func readCsv(filepath, delimiter string) ([]CarData, error) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dateStr := strings.Split(path.Base(filepath), " ")[0]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(date)

	FullData := []CarData{}

	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return FullData, err
	}
	first := true
	for scanner.Scan() {
		if first {
			first = false
			continue
		}
		data := lineToData(scanner.Text(), delimiter, date)
		FullData = append(FullData, data)
	}

	return FullData, err
}

func lineToData(line, delimiter string, date time.Time) CarData {
	data := CarData{}
	array := strings.Split(line, delimiter)
	if s, err := strconv.ParseFloat(strings.Trim(array[0], "\""), 64); err == nil {
		data.Second = s
	} else {
		log.Printf("Unable to parse value %s as float", array[0])
	}
	data.Pid = strings.Trim(array[1], "\"")
	if s, err := strconv.ParseFloat(strings.Trim(array[2], "\""), 64); err == nil {
		data.Value = s
	} else {
		log.Printf("Unable to parse value %s as float", array[2])
	}
	data.Units = strings.Trim(array[3], "\"")

	return data
}
