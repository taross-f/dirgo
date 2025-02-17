package output

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/taross-f/dirgo/pkg/analyzer"
)

func TestFormatter_WriteResults(t *testing.T) {
	results := []analyzer.Result{
		{
			Path:  "/path/to/dir1",
			Size:  1024 * 1024, // 1MB
			Count: 10,
		},
		{
			Path:  "/path/to/dir2",
			Size:  512 * 1024, // 512KB
			Count: 5,
		},
	}

	tests := []struct {
		name     string
		results  []analyzer.Result
		limit    int
		wantRows int
	}{
		{
			name:     "full output",
			results:  results,
			limit:    2,
			wantRows: 2,
		},
		{
			name:     "limited output",
			results:  results,
			limit:    1,
			wantRows: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := NewFormatter(&buf)

			err := f.WriteResults(tt.results, tt.limit)
			if err != nil {
				t.Fatalf("WriteResults() error = %v", err)
			}

			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			if len(lines) != tt.wantRows {
				t.Errorf("WriteResults() got %d rows, want %d", len(lines), tt.wantRows)
			}
		})
	}
}

func TestFormatter_WriteCSV(t *testing.T) {
	results := []analyzer.Result{
		{
			Path:  "/path/to/dir1",
			Size:  1024 * 1024,
			Count: 10,
		},
		{
			Path:  "/path/to/dir2",
			Size:  512 * 1024,
			Count: 5,
		},
	}

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "test.csv")

	f := NewFormatter(os.Stdout)
	err := f.WriteCSV(results, csvPath)
	if err != nil {
		t.Fatalf("WriteCSV() error = %v", err)
	}

	// ファイルが作成されたことを確認
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		t.Error("WriteCSV() did not create the file")
	}

	// ファイルの内容を読み込んで検証
	content, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 3 { // ヘッダー + 2行のデータ
		t.Errorf("WriteCSV() got %d rows, want 3", len(lines))
	}

	// ヘッダーの検証
	if !strings.HasPrefix(lines[0], "Path,Size,FileCount") {
		t.Error("WriteCSV() header is incorrect")
	}
} 