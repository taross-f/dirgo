package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/taross-f/dirgo/pkg/analyzer"
	"github.com/taross-f/dirgo/pkg/output"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "dirgo",
		Usage:     "ディレクトリサイズ解析ツール",
		UsageText: "dirgo [options] <target_path>",
		Version:   "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "outfile",
				Aliases: []string{"o"},
				Usage:   "出力CSVファイルのパス",
			},
			&cli.IntFlag{
				Name:    "async-depth",
				Aliases: []string{"d"},
				Value:   3,
				Usage:   "非同期処理を行う深さ",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "詳細なログを出力",
			},
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"c"},
				Value:   20,
				Usage:   "出力する結果の数",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("ターゲットパスを指定してください")
	}

	// 設定の取得
	targetPath := c.Args().Get(0)
	outputPath := c.String("outfile")
	asyncDepth := c.Int("async-depth")
	verbose := c.Bool("verbose")
	count := c.Int("count")

	if verbose {
		log.Printf("Target path: %s", targetPath)
		log.Printf("Output path: %s", outputPath)
		log.Printf("Async depth: %d", asyncDepth)
		log.Printf("CPU cores: %d", runtime.NumCPU())
	}

	// アナライザーの初期化
	a := analyzer.NewAnalyzer(1, asyncDepth, runtime.NumCPU())

	// 解析の実行
	results, err := a.Analyze(targetPath)
	if err != nil {
		return fmt.Errorf("解析エラー: %w", err)
	}

	if verbose {
		log.Printf("Found %d directories", len(results))
	}

	// 結果の出力
	f := output.NewFormatter(os.Stdout)
	if err := f.WriteResults(results, count); err != nil {
		return fmt.Errorf("結果の出力エラー: %w", err)
	}

	// CSVファイルへの出力
	if outputPath != "" {
		if err := f.WriteCSV(results, outputPath); err != nil {
			return fmt.Errorf("CSVファイルの出力エラー: %w", err)
		}
		if verbose {
			log.Printf("Results written to %s", outputPath)
		}
	}

	return nil
}
