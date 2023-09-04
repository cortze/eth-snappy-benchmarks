package metrics

import (
	"fmt"
	"time"

	"github.com/cortze/eth-snappy-benchmarks/csvs"
)

type aggregator string

func (a aggregator) String() string { return string(a) }

type metric string

func (m metric) String() string { return string(m) }

const (
	Minimum aggregator = "min"
	Maximum aggregator = "max"
	Average aggregator = "avg"

	FolderName    metric = "folder"
	FileName      metric = "file"
	RawSize       metric = "raw-size"
	CompressSize  metric = "compress-size"
	EncodingTime  metric = "encoding-time"
	DecodingTime  metric = "decoding-time"
	CompressRatio metric = "compress-ratio"
	CompressSpeed metric = "compress-speed"
)

type CompressionMetrics struct {
	// benchmark details
	BlockFolder   string
	MetricsFolder string
	MetricsFile   string
	FileNames     []string
	RawSize       []int64
	CompressSize  []int64
	// benchmark data
	CompressRatios []float64
	CompressSpeeds []float64
	EncodeSpeeds   []time.Duration // (Bytes/Millisecond)
	DecodeSpeeds   []time.Duration // (Bytes/Millisecond)
	// csv related
	csvExporter *csvs.CSV
	csvColumns  []csvs.Stringable
}

func NewCompressionMetrics(blocksFolder, metricsFolder, metricsFile, encoding string) (*CompressionMetrics, error) {
	csvColumns := []csvs.Stringable{
		FolderName, FileName, RawSize, CompressSize,
		EncodingTime, DecodingTime, CompressRatio, CompressSpeed}

	csvFile, err := csvs.NewCsvExporter(metricsFolder+"/"+encoding+"_"+metricsFile, csvColumns)
	if err != nil {
		return nil, err
	}

	return &CompressionMetrics{
		BlockFolder:   blocksFolder,
		MetricsFolder: metricsFolder,
		MetricsFile:   metricsFile,
		FileNames:     make([]string, 0),
		RawSize:       make([]int64, 0),
		CompressSize:  make([]int64, 0),
		// metrics
		EncodeSpeeds:   make([]time.Duration, 0),
		DecodeSpeeds:   make([]time.Duration, 0),
		CompressRatios: make([]float64, 0),
		CompressSpeeds: make([]float64, 0),
		//
		csvExporter: csvFile,
		csvColumns:  csvColumns}, nil
}

func (m *CompressionMetrics) AddResults(
	file string, rawSize, compressSize int64,
	encodingSpeed, decodingSpeed time.Duration,
	compressRatio, compressSpeed float64) {

	// file related
	m.FileNames = append(m.FileNames, file)
	m.RawSize = append(m.RawSize, rawSize)
	m.CompressSize = append(m.CompressSize, compressSize)
	// compression related
	m.EncodeSpeeds = append(m.EncodeSpeeds, encodingSpeed)
	m.DecodeSpeeds = append(m.DecodeSpeeds, decodingSpeed)
	m.CompressRatios = append(m.CompressRatios, compressRatio)
	m.CompressSpeeds = append(m.CompressSpeeds, compressSpeed)
}

func (m *CompressionMetrics) GetSummary() map[aggregator]map[metric]float64 {
	var avgEncode time.Duration = 0
	var avgDecode time.Duration = 0
	var avgRatio float64 = 0
	var avgSpeed float64 = 0

	items := len(m.FileNames)
	for i := 0; i < items; i++ {
		avgRatio = avgRatio + m.CompressRatios[i]
		avgSpeed = avgSpeed + m.CompressSpeeds[i]
		avgEncode = avgEncode + m.EncodeSpeeds[i]
		avgDecode = avgDecode + m.DecodeSpeeds[i]
	}
	encodeMin, encodeMax := findMinAndMax(m.EncodeSpeeds)
	decodeMin, decodeMax := findMinAndMax(m.DecodeSpeeds)
	ratioMin, ratioMax := findMinAndMax(m.CompressRatios)
	speedMin, speedMax := findMinAndMax(m.CompressSpeeds)

	summary := make(map[aggregator]map[metric]float64, 3)

	minSummary := map[metric]float64{
		EncodingTime:  float64(encodeMin.Nanoseconds()),
		DecodingTime:  float64(decodeMin.Nanoseconds()),
		CompressRatio: ratioMin,
		CompressSpeed: speedMin,
	}

	maxSummary := map[metric]float64{
		EncodingTime:  float64(encodeMax.Nanoseconds()),
		DecodingTime:  float64(decodeMax.Nanoseconds()),
		CompressRatio: ratioMax,
		CompressSpeed: speedMax,
	}

	avgSummary := map[metric]float64{
		EncodingTime:  float64(avgEncode.Nanoseconds()) / float64(items),
		DecodingTime:  float64(avgEncode.Nanoseconds()) / float64(items),
		CompressRatio: float64(avgRatio) / float64(items),
		CompressSpeed: float64(avgSpeed) / float64(items),
	}

	summary[Minimum] = minSummary
	summary[Maximum] = maxSummary
	summary[Average] = avgSummary
	return summary
}

func findMinAndMax[T time.Duration | float64](a []T) (min T, max T) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func (m *CompressionMetrics) rowComposer(rawRow []interface{}) []string {
	row := make([]string, len(rawRow), len(rawRow))
	for idx, item := range rawRow {
		switch item.(type) {
		case float64:
			row[idx] = fmt.Sprintf("%.6f", item.(float64))
		case int64:
			row[idx] = fmt.Sprintf("%d", item.(int64))
		case string:
			row[idx] = item.(string)
		case time.Duration:
			newItem := item.(time.Duration)
			row[idx] = fmt.Sprintf("%.6f", float64(newItem.Nanoseconds()))
		default:
			row[idx] = fmt.Sprint(item)
		}
	}
	return row
}

func (m *CompressionMetrics) Export() error {
	rows := m.generateRows()
	return m.csvExporter.Export(rows, m.rowComposer)
}

func (m *CompressionMetrics) generateRows() [][]interface{} {
	rows := make([][]interface{}, 0)
	// combine the arrays into rows for the csv
	for idx, fileName := range m.FileNames {
		row := make([]interface{}, 0)
		row = append(row, m.BlockFolder)
		row = append(row, fileName)
		row = append(row, m.RawSize[idx])
		row = append(row, m.CompressSize[idx])
		row = append(row, m.EncodeSpeeds[idx])
		row = append(row, m.DecodeSpeeds[idx])
		row = append(row, m.CompressRatios[idx])
		row = append(row, m.CompressSpeeds[idx])

		rows = append(rows, row)
	}
	return rows
}
