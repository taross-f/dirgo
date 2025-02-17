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
			wantDirs: 2,
			wantFiles: map[string]uint64{
				"dir1": 2,
				"dir2": 1,
			},
			wantSize: map[string]uint64{
				"dir1": 300, // 100 + 200
				"dir2": 300,
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
