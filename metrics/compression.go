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
	CompressRatios []time.Duration
	CompressSpeeds []time.Duration
	EncodeSpeeds   []float64 // (Bytes/Millisecond)
	DecodeSpeeds   []float64 // (Bytes/Millisecond)
	// csv related
	csvExporter *csvs.CSV
	csvColumns  []csvs.Stringable
}

func NewCompressionMetrics(blocksFolder, metricsFolder, metricsFile string) (*CompressionMetrics, error) {
	csvColumns := []csvs.Stringable{
		FolderName, FileName, RawSize, CompressSize,
		EncodingTime, DecodingTime, CompressRatio, CompressSpeed}

	csvFile, err := csvs.NewCsv(metricsFolder+"/"+metricsFile, csvColumns)
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
		EncodeSpeeds:   make([]float64, 0),
		DecodeSpeeds:   make([]float64, 0),
		CompressRatios: make([]time.Duration, 0),
		CompressSpeeds: make([]time.Duration, 0),
		//
		csvExporter: csvFile,
		csvColumns:  csvColumns}, nil
}

func (m *CompressionMetrics) AddResults(
	file string, rawSize, compressSize int64,
	compressRatio, compressSpeed time.Duration,
	encodingSpeed, decodingSpeed float64) {

	// file related
	m.FileNames = append(m.FileNames, file)
	m.RawSize = append(m.RawSize, rawSize)
	m.CompressSize = append(m.CompressSize, compressSize)
	// compression related
	m.CompressRatios = append(m.CompressRatios, compressRatio)
	m.CompressSpeeds = append(m.CompressRatios, compressSpeed)
	m.EncodeSpeeds = append(m.EncodeSpeeds, encodingSpeed)
	m.DecodeSpeeds = append(m.DecodeSpeeds, decodingSpeed)
}

func (m *CompressionMetrics) GetSummary(target time.Duration) map[aggregator]map[metric]float64 {
	durationConversion := target
	floatConversion := float64(time.Nanosecond / durationConversion)

	var avgEncode float64 = 0
	var avgDecode float64 = 0
	var avgRatio time.Duration = 0
	var avgSpeed time.Duration = 0

	items := len(m.FileNames)
	for i := 0; i < items; i++ {
		avgRatio = avgRatio + m.CompressRatios[i]
		avgSpeed = avgSpeed + m.CompressSpeeds[i]
		avgEncode = avgEncode + m.EncodeSpeeds[i]
		avgDecode = avgDecode + m.DecodeSpeeds[i]
	}
	EncodeMin, EncodeMax := findMinAndMax(m.EncodeSpeeds)
	DecodeMin, DecodeMax := findMinAndMax(m.DecodeSpeeds)
	ratioMin, ratioMax := findMinAndMax(m.CompressRatios)
	speedMin, speedMax := findMinAndMax(m.CompressSpeeds)

	summary := make(map[aggregator]map[metric]float64, 3)

	minSummary := map[metric]float64{
		EncodingTime:  EncodeMin,
		DecodingTime:  DecodeMin,
		CompressRatio: float64(ratioMin) / floatConversion,
		CompressSpeed: float64(speedMin) / floatConversion,
	}

	maxSummary := map[metric]float64{
		EncodingTime:  EncodeMax,
		DecodingTime:  DecodeMax,
		CompressRatio: float64(ratioMax) / floatConversion,
		CompressSpeed: float64(speedMax) / floatConversion,
	}

	avgSummary := map[metric]float64{
		EncodingTime:  avgEncode / float64(items),
		DecodingTime:  avgEncode / float64(items),
		CompressRatio: (float64(avgRatio) / floatConversion) / float64(items),
		CompressSpeed: (float64(avgSpeed) / floatConversion) / float64(items),
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
