package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzer_Analyze(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用のディレクトリ構造を作成
	dirs := []string{
		"dir1",
		"dir1/subdir1",
		"dir2",
		"empty_dir",
	}

	files := map[string]int64{
		"dir1/file1.txt":         100,
		"dir1/subdir1/file2.txt": 200,
		"dir2/file3.txt":         300,
	}

	// ディレクトリを作成
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("failed to create test directory: %v", err)
		}
	}

	// ファイルを作成
	for path, size := range files {
		content := make([]byte, size)
		err := os.WriteFile(filepath.Join(tempDir, path), content, 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// テストケース
	tests := []struct {
		name      string
		maxDepth  int
		wantDirs  int
		wantFiles map[string]uint64 // パスごとのファイル数
		wantSize  map[string]uint64 // パスごとのサイズ
	}{
		{
			name:     "depth 1",
			maxDepth: 1,
			wantDirs: 4, // dir1, dir2, empty_dir
			wantFiles: map[string]uint64{
				"dir1":      1, // file1.txt
				"dir2":      1, // file3.txt
				"empty_dir": 0,
			},
			wantSize: map[string]uint64{
				"dir1":      100, // file1.txt
				"dir2":      300, // file3.txt
				"empty_dir": 0,
			},
		},
		{
			name:     "depth 2",
			maxDepth: 2,
			wantDirs: 4, // dir1, dir1/subdir1, dir2, empty_dir
			wantFiles: map[string]uint64{
				"dir1":      1, // file1.txt
				"subdir1":   1, // file2.txt
				"dir2":      1, // file3.txt
				"empty_dir": 0,
			},
			wantSize: map[string]uint64{
				"dir1":      100, // file1.txt
				"subdir1":   200, // file2.txt
				"dir2":      300, // file3.txt
				"empty_dir": 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnalyzer(tt.maxDepth, 1, 1)
			results, err := a.Analyze(tempDir)
			if err != nil {
				t.Fatalf("Analyze() error = %v", err)
			}

			if len(results) != tt.wantDirs {
				t.Errorf("Analyze() got %d directories, want %d", len(results), tt.wantDirs)
			}

			// 結果を検証
			for _, result := range results {
				basePath := filepath.Base(result.Path)
				if wantFiles, ok := tt.wantFiles[basePath]; ok {
					if result.Count != wantFiles {
						t.Errorf("Analyze() for %s got %d files, want %d", basePath, result.Count, wantFiles)
					}
				}

				if wantSize, ok := tt.wantSize[basePath]; ok {
					if result.Size != wantSize {
						t.Errorf("Analyze() for %s got size %d, want %d", basePath, result.Size, wantSize)
					}
				}
			}
		})
	}
}

func TestAnalyzer_AnalyzeErrors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "non-existent directory",
			path:    "/path/that/does/not/exist",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnalyzer(1, 1, 1)
			_, err := a.Analyze(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResult_String(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		want   string
	}{
		{
			name: "normal result",
			result: Result{
				Path:  "/test/path",
				Size:  1024,
				Count: 5,
			},
			want: "/test/path: 1024 bytes, 5 files",
		},
		{
			name: "empty result",
			result: Result{
				Path:  "/empty",
				Size:  0,
				Count: 0,
			},
			want: "/empty: 0 bytes, 0 files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.String(); got != tt.want {
				t.Errorf("Result.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
