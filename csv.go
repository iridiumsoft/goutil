package util

import (
	"os"
	"io"
	"log"
	"strings"
	"encoding/csv"
	"github.com/sirupsen/logrus"
)

func CSV2JSON(Path string) []Object {
	clientsFile, _ := os.OpenFile(Path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	r := csv.NewReader(clientsFile)
	var columns []string
	var data []Object
	i := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("oh no")
		}
		if err != nil {
			logrus.Error(err.Error())
		}
		i++

		if i == 1 { // Get Columns
			columns = record
			continue
		}

		if record[0] == "N/A" {
			continue
		}
		// Make Key pair Object
		row := make(Object)
		for inc, item := range record {
			// Convert keys to smaller case and
			Key := strings.ToLower(columns[inc])
			Key = strings.Replace(Key, " ", "", -1)
			if row[Key] != nil {
				Key = Key + "_old"
			}
			if Key == "" {
				delete(row, "")
				continue
			}
			if Key == "Villages" || Key == "villages" {
				Key = "village"
			}
			if Key == "na" {
				item = strings.Replace(item, "NA-", "", -1)
			}
			row[Key] = item
		}
		data = append(data, row)
	}
	return data
}
