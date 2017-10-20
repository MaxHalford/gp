package dataset

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadCSV(path string, target string, classification bool) (*Dataset, error) {
	var (
		dataset = &Dataset{}
		f, _    = os.Open(path)
		r       = csv.NewReader(bufio.NewReader(f))
	)
	defer f.Close()

	// Read the headers
	columns, err := r.Read()
	if err != nil {
		return nil, err
	}

	// Go through column names and determine which one is the target
	var (
		targetIdx = -1
		p         int
	)
	dataset.XNames = make([]string, len(columns)-1)
	for i, column := range columns {
		if column == target {
			dataset.YName = column
			targetIdx = i
			break
		} else {
			if p == len(dataset.XNames) {
				break
			}
			dataset.XNames[p] = column
			p++
		}
	}
	if targetIdx == -1 {
		return nil, fmt.Errorf("No column named '%s'", target)
	}

	// Initialize the columns
	dataset.X = make([][]float64, p)

	// Initialize an empty class map in case of classification
	if classification {
		dataset.ClassMap = MakeClassMap()
	}

	// Iterate over the rows
	for {
		var record, err = r.Read()
		if err == io.EOF {
			break
		}
		// Parse the features as float64s
		for i, s := range record {
			if i != targetIdx {
				x, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return nil, err
				}
				dataset.X[i] = append(dataset.X[i], x)
			}
		}
		// Parse the target
		if classification {
			dataset.Y = append(dataset.Y, dataset.ClassMap.Get(record[targetIdx]))
		} else {
			var y, _ = strconv.ParseFloat(record[targetIdx], 64)
			dataset.Y = append(dataset.Y, y)
		}
	}
	return dataset, nil
}
