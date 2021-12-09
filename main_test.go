package carscannertodb

import (
	"reflect"
	"testing"
	"time"
)

func Test_lineToData(t *testing.T) {
	type args struct {
		line      string
		delimiter string
	}
	tests := []struct {
		name string
		args args
		want CarData
	}{
		// TODO: Add test cases.
	}
	tests = append(tests,
		struct {
			name string
			args args
			want CarData
		}{
			name: "simple",
			args: args{line: "\"60338.071981\";\"Обороты двигателя\";\"994\";\"rpm\"", delimiter: ";"},
			want: CarData{},
		},
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lineToData(tt.args.line, tt.args.delimiter, time.Now()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lineToData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadCsv(t *testing.T) {
	readCsv("/Users/roe/CarLogs/CarScanner/2021-11-19 16-45-25.csv", ";")
}
