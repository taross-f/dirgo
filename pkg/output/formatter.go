package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/taross-f/dirgo/pkg/analyzer"
)

// Formatter は解析結果の出力を整形する構造体です
type Formatter struct {
	writer io.Writer
}

// NewFormatter は新しいFormatterインスタンスを作成します
func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{writer: w}
}

// WriteResults は解析結果を指定されたフォーマットで出力します
func (f *Formatter) WriteResults(results []analyzer.Result, limit int) error {
	// サイズでソート
	sort.Slice(results, func(i, j int) bool {
		return results[i].Size > results[j].Size
	})

	// 上位n件を出力
	for i := 0; i < limit && i < len(results); i++ {
		fmt.Fprintf(f.writer, "%s,%s,%d\n",
			results[i].Path,
			humanize.Bytes(results[i].Size),
			results[i].Count,
		)
	}

	return nil
}

// WriteCSV は解析結果をCSVファイルに出力します
func (f *Formatter) WriteCSV(results []analyzer.Result, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込み
	if err := writer.Write([]string{"Path", "Size", "FileCount"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// データを書き込み
	for _, result := range results {
		if err := writer.Write([]string{
			result.Path,
			fmt.Sprint(result.Size),
			fmt.Sprint(result.Count),
		}); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
