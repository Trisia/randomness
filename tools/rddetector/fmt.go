package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// ReportFormatter 报告格式化接口
type ReportFormatter interface {
	// FormatTestReport 格式化检测报告
	// - results: 检测结果
	// - w: 输出流
	// - 返回值: 错误信息
	FormatTestReport(results []*R, w io.Writer) error
	// FormatAnalysisReport 格式化分析报告
	// - results: 检测结果
	// - w: 输出流
	// - 返回值: 错误信息
	FormatAnalysisReport(results []AnalysisResult, w io.Writer) error
}

// ReportCollector 统一数据收集器
type ReportCollector struct {
	results      []*R       // 检测结果
	mu           sync.Mutex // 互斥锁
	format       string     // 输出格式 (csv/json)
	reportPath   string     // 报告输出路径
	analysisPath string     // 分析报告输出路径
	threshold    float64    // 通过判定阈值
}

// NewReportCollector 创建新的数据收集器
// - format: 输出格式 (csv/json/xml)
// - reportPath: 报告输出路径
// - analysisPath: 分析报告输出路径
// - threshold: 通过判定阈值 (0.0-1.0)
func NewReportCollector(format, reportPath, analysisPath string, threshold float64) *ReportCollector {
	return &ReportCollector{
		results:      make([]*R, 0),
		format:       format,
		reportPath:   reportPath,
		analysisPath: analysisPath,
		threshold:    threshold,
	}
}

// AddResult 添加检测结果
func (c *ReportCollector) AddResult(result *R) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results = append(c.results, result)
}

// GetResults 获取所有检测结果
func (c *ReportCollector) GetResults() []*R {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.results
}

// GenerateReports 生成检测报告和分析报告
func (c *ReportCollector) GenerateReports() error {
	results := c.GetResults()

	// 生成检测报告
	if c.reportPath != "" && len(results) > 0 {
		if err := c.generateTestReport(results, c.reportPath); err != nil {
			return fmt.Errorf("生成检测报告失败: %v", err)
		}
	}

	// 生成分析报告
	if c.analysisPath != "" && len(results) > 0 {
		if err := c.generateAnalysisReport(results, c.analysisPath); err != nil {
			return fmt.Errorf("生成分析报告失败: %v", err)
		}
	}

	return nil
}

// generateTestReport 生成检测报告
func (c *ReportCollector) generateTestReport(results []*R, outputPath string) error {
	// 创建输出文件
	_ = os.MkdirAll(filepath.Dir(outputPath), os.FileMode(0600))
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer file.Close()

	formatter := getFormatter(c.format)
	return formatter.FormatTestReport(results, file)
}

// generateAnalysisReport 生成分析报告
func (c *ReportCollector) generateAnalysisReport(results []*R, outputPath string) error {
	// 统计每个检测项目的通过情况
	var analysisResults []AnalysisResult
	totalFiles := len(results)

	if totalFiles == 0 {
		return nil
	}

	// 遍历所有检测项目
	for i, testItem := range results[0].TestItems {
		passCount := 0

		for _, result := range results {
			if i < len(result.TestItems) {
				// 根据GM/T 0005-2021规范：P >= 0.01 且 Q >= 0.0001
				if result.TestItems[i].PValue >= 0.01 && result.TestItems[i].QValue >= 0.0001 {
					passCount++
				}
			}
		}

		passRate := float64(passCount) / float64(totalFiles)
		isPassed := passRate >= c.threshold

		analysisResults = append(analysisResults, AnalysisResult{
			TestName:    testItem.TestName,
			PassCount:   passCount,
			TotalCount:  totalFiles,
			PassRate:    passRate,
			Requirement: c.threshold,
			IsPassed:    isPassed,
		})
	}

	// 创建输出文件
	_ = os.MkdirAll(filepath.Dir(outputPath), os.FileMode(0600))
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer file.Close()

	formatter := getFormatter(c.format)
	return formatter.FormatAnalysisReport(analysisResults, file)
}

// JSONFormatter JSON格式输出
type JSONFormatter struct{}

func (f *JSONFormatter) FormatTestReport(results []*R, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func (f *JSONFormatter) FormatAnalysisReport(results []AnalysisResult, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

// CSVFormatter CSV格式输出
type CSVFormatter struct{}

func (f *CSVFormatter) FormatTestReport(results []*R, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 写入CSV表头
	headers := []string{"文件名"}
	if len(results) > 0 {
		for _, item := range results[0].TestItems {
			headers = append(headers, item.TestName+" P值", item.TestName+" Q值")
		}
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 写入数据
	for _, result := range results {
		record := []string{result.Name}
		for _, item := range result.TestItems {
			record = append(record,
				fmt.Sprintf("%.6f", item.PValue),
				fmt.Sprintf("%.6f", item.QValue))
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func (f *CSVFormatter) FormatAnalysisReport(results []AnalysisResult, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 写入CSV表头
	headers := []string{"检测项目（含参数）", "通过数", "检测数", "通过率", "满足随机性要求", "是否通过"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 写入数据
	for _, result := range results {
		passRateStr := fmt.Sprintf("%.4f", result.PassRate)
		requirementStr := fmt.Sprintf("%.3f", result.Requirement)
		isPassedStr := "是"
		if !result.IsPassed {
			isPassedStr = "否"
		}

		record := []string{
			result.TestName,
			strconv.Itoa(result.PassCount),
			strconv.Itoa(result.TotalCount),
			passRateStr,
			requirementStr,
			isPassedStr,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// XMLFormatter XML格式输出
type XMLFormatter struct{}

func (f *XMLFormatter) FormatTestReport(results []*R, w io.Writer) error {
	_, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<RandomnessTestReport>\n"))
	if err != nil {
		return err
	}

	for _, result := range results {
		_, err = w.Write([]byte(fmt.Sprintf("  <File name=\"%s\">\n", result.Name)))
		if err != nil {
			return err
		}

		for _, item := range result.TestItems {
			_, err = w.Write([]byte(fmt.Sprintf("    <Test name=\"%s\" p=\"%.6f\" q=\"%.6f\"/>\n",
				item.TestName, item.PValue, item.QValue)))
			if err != nil {
				return err
			}
		}

		_, err = w.Write([]byte("  </File>\n"))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("</RandomnessTestReport>\n"))
	return err
}

func (f *XMLFormatter) FormatAnalysisReport(results []AnalysisResult, w io.Writer) error {
	_, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<AnalysisReport>\n"))
	if err != nil {
		return err
	}

	for _, result := range results {
		isPassedStr := "true"
		if !result.IsPassed {
			isPassedStr = "false"
		}

		_, err = w.Write([]byte(fmt.Sprintf("  <Test name=\"%s\" passCount=\"%d\" totalCount=\"%d\" passRate=\"%.4f\" requirement=\"%.3f\" isPassed=\"%s\"/>\n",
			result.TestName, result.PassCount, result.TotalCount, result.PassRate, result.Requirement, isPassedStr)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("</AnalysisReport>\n"))
	return err
}

// getFormatter 根据格式名称获取格式化器
func getFormatter(format string) ReportFormatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "xml":
		return &XMLFormatter{}
	case "csv":
		fallthrough
	default:
		return &CSVFormatter{}
	}
}
