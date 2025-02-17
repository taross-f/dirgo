# dirgo

ディレクトリサイズを効率的に解析するGoツール

## 機能

- ディレクトリのサイズと含まれるファイル数を解析
- 非同期処理による高速な解析
- 結果のCSV出力
- 人間が読みやすい形式でのサイズ表示

## インストール

```bash
go install github.com/taross-f/dirgo/cmd/dirgo@latest
```

## 使用方法

基本的な使用方法：

```bash
dirgo [オプション] <対象ディレクトリ>
```

### オプション

- `-o, --outfile`: 結果をCSVファイルに出力
- `-d, --async-depth`: 非同期処理を行う深さ（デフォルト: 3）
- `-v, --verbose`: 詳細なログを出力
- `-c, --count`: 出力する結果の数（デフォルト: 20）

### 使用例

```bash
# 基本的な使用方法
dirgo /path/to/directory

# 詳細なログ出力を有効にする
dirgo -v /path/to/directory

# 結果をCSVファイルに出力
dirgo -o results.csv /path/to/directory

# 表示する結果の数を指定
dirgo -c 10 /path/to/directory

# 非同期処理の深さを指定
dirgo -d 4 /path/to/directory
```

## 出力形式

### 標準出力
```
/path/to/dir1,1.2GB,150
/path/to/dir2,800MB,42
```

### CSVファイル
```csv
Path,Size,FileCount
/path/to/dir1,1288490189,150
/path/to/dir2,839120938,42
```

## ビルド

```bash
git clone https://github.com/taross-f/dirgo.git
cd dirgo
go build -o dirgo cmd/dirgo/main.go
```

## 開発

### 必要条件

- Go 1.22以上
- golangci-lint（開発時）

### テストの実行

```bash
go test -v ./...
```

### リンターの実行

```bash
golangci-lint run
```

## ライセンス

MIT License

Copyright (c) 2024 taross-f

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.



