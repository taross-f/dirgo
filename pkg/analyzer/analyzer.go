package analyzer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// Result は解析結果を表す構造体です
type Result struct {
	Path  string
	Size  uint64
	Count uint64
}

// String はResult構造体の文字列表現を返します
func (r *Result) String() string {
	return fmt.Sprintf("%s: %d bytes, %d files", r.Path, r.Size, r.Count)
}

// Analyzer はディレクトリ解析を行う構造体です
type Analyzer struct {
	maxDepth   int
	asyncDepth int
	workers    int
}

// NewAnalyzer は新しいAnalyzerインスタンスを作成します
func NewAnalyzer(maxDepth, asyncDepth, workers int) *Analyzer {
	return &Analyzer{
		maxDepth:   maxDepth,
		asyncDepth: asyncDepth,
		workers:    workers,
	}
}

// Analyze は指定されたルートディレクトリの解析を行います
func (a *Analyzer) Analyze(root string) ([]Result, error) {
	if root == "" {
		return nil, fmt.Errorf("empty path provided")
	}

	root = filepath.Clean(root)
	if _, err := os.Stat(root); err != nil {
		return nil, fmt.Errorf("invalid path %s: %w", root, err)
	}

	paths := a.getTargetPaths(root, 0)
	results := make([]Result, 0, len(paths))

	var wg sync.WaitGroup
	resultChan := make(chan Result, a.workers)
	doneChan := make(chan struct{})

	// 結果を収集するゴルーチン
	go func() {
		for result := range resultChan {
			results = append(results, result)
		}
		close(doneChan)
	}()

	// 各パスを非同期で解析
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			size, count := a.analyzeDir(p)
			resultChan <- Result{
				Path:  p,
				Size:  size,
				Count: count,
			}
		}(path)
	}

	wg.Wait()
	close(resultChan)
	<-doneChan // 結果の収集が完了するまで待機

	return results, nil
}

// getTargetPaths は解析対象のパスリストを取得します
func (a *Analyzer) getTargetPaths(root string, depth int) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}

	paths := make([]string, 0)
	if depth >= a.maxDepth {
		for _, entry := range entries {
			if entry.IsDir() {
				paths = append(paths, filepath.Join(root, entry.Name()))
			}
		}
		return paths
	}

	// 現在のディレクトリのサブディレクトリを追加
	for _, entry := range entries {
		if entry.IsDir() {
			subPath := filepath.Join(root, entry.Name())
			paths = append(paths, subPath)
		}
	}

	// サブディレクトリの下のパスを追加
	for _, entry := range entries {
		if entry.IsDir() {
			subPath := filepath.Join(root, entry.Name())
			subPaths := a.getTargetPaths(subPath, depth+1)
			paths = append(paths, subPaths...)
		}
	}

	return paths
}

// analyzeDir はディレクトリのサイズと含まれるファイル数を計算します
func (a *Analyzer) analyzeDir(path string) (uint64, uint64) {
	var size, count uint64

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		// ルートディレクトリ以外のディレクトリはスキップ
		if d.IsDir() && p != path {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return nil
			}
			size += uint64(info.Size())
			count++
		}
		return nil
	})

	if err != nil {
		return 0, 0
	}

	return size, count
}
