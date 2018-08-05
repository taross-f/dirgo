package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/urfave/cli"
)

const maxDepth int = 1

var asyncDepth int
var outputPath string
var vervose bool
var count int

type Output struct {
	Path  string
	Size  uint64
	Count uint64
}

func (o *Output) str() string {
	return o.Path + "," + fmt.Sprint(humanize.Bytes(o.Size)) + "," + fmt.Sprint(o.Count)
}

type BySize []Output

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size > a[j].Size } // descendent

var cpuCount int

func main() {
	app := cli.NewApp()
	app.Name = "dirgo"
	app.HelpName = "dirgo"
	app.UsageText = "dirgo [-o output_path] [-d asyncDepth]  target_path"
	app.Version = "0.0.1"
	app.Action = core
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "outfile, o",
			Usage:       "file path to output",
			Destination: &outputPath,
		},
		cli.IntFlag{
			Name:        "asyncdepth, d",
			Usage:       "depth searching asynchronously",
			Value:       3,
			Destination: &asyncDepth,
		},
		cli.BoolFlag{
			Name:        "vervose, V",
			Usage:       "log vervosely",
			Destination: &vervose,
		},
		cli.IntFlag{
			Name:        "count, c",
			Usage:       "output path count",
			Value:       20,
			Destination: &count,
		},
	}
	app.Run(os.Args)
}

func log(s string) {
	if !vervose {
		return
	}
	println(s)
}

func core(c *cli.Context) error {
	cpuCount = runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	log("start! Core:" + string(cpuCount))
	root := c.Args().Get(0)
	_, err := os.Stat(root)
	if err != nil {
		fmt.Println("You must set the valid target path.")
		return err
	}
	root = strings.TrimRight(root, "/\\")
	log("output: " + outputPath)
	log("async depth: " + string(asyncDepth))

	checkNonRepeat(root)
	return nil
}

func checkNonRepeat(root string) {
	syncStart := time.Now()
	paths := getTargetPaths(root, 0)
	log("path count:" + fmt.Sprint(len(paths)))
	buf := ""

	result := make(chan Output, cpuCount)
	c := make(chan Output)
	forSort := []Output{}
	go getSizeRecursiveNonRepeat(root, 0, c, result)
	for i := 0; i < len(paths); i++ {
		o := <-result
		forSort = append(forSort, o)
		buf += o.Path + "," + fmt.Sprint(o.Size) + "," + fmt.Sprint(o.Count) + "\n"
	}

	sort.Sort(BySize(forSort))
	for i := 0; i < count; i++ {
		fmt.Println(forSort[i].str())
	}

	if outputPath != "" {
		ioutil.WriteFile(outputPath, []byte(buf), os.ModePerm)
	}

	syncEnd := time.Now()

	log(fmt.Sprintf("---output: %f sec", syncEnd.Sub(syncStart).Seconds()))
}

func getTargetPaths(root string, depth int) []string {
	fi, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Println("error occured: ", err.Error())
		return make([]string, 0) // if permission denied, return empty
	}

	paths := make([]string, 0)
	if depth >= maxDepth {
		for _, f := range fi {
			if f.IsDir() {
				paths = append(paths, root+"/"+f.Name())
			}
		}
		return paths
	}

	for _, f := range fi {
		if f.IsDir() {
			paths = append(paths, getTargetPaths(root+"/"+f.Name(), depth+1)...)
			paths = append(paths, root+"/"+f.Name())
		}
	}
	return paths
}

func getSizeRecursive(root, search string) (uint64, uint64) {
	fi, err := ioutil.ReadDir(search)
	if err != nil {
		return 0, 0 // if permission denied, return zeros
	}

	var size, count uint64
	for _, f := range fi {
		if f.IsDir() {
			n := f.Name()
			s, c := getSizeRecursive(root, search+"/"+n)
			size += s
			count += c
		} else {
			size += uint64(f.Size())
			count++
		}
	}
	return size, count
}

func getSizeRecursiveNonRepeat(search string, depth int, outputChan chan Output, resultChan chan Output) {
	fi, err := ioutil.ReadDir(search)
	if err != nil {
		outputChan <- Output{Path: search, Size: 0, Count: 0}
		if depth <= maxDepth+1 {
			resultChan <- Output{Path: search, Size: 0, Count: 0}
		}
		return
	}

	var size, count uint64
	length := 0
	for _, f := range fi {
		if f.IsDir() {
			length++
		}
	}
	nextOutput := make(chan Output)
	defer close(nextOutput)
	for _, f := range fi {
		if f.IsDir() {
			if depth <= asyncDepth {
				go getSizeRecursiveNonRepeat(search+string(os.PathSeparator)+f.Name(), depth+1, nextOutput, resultChan)
			} else {
				s, c := getSizeRecursive(search+string(os.PathSeparator)+f.Name(), search+string(os.PathSeparator)+f.Name())
				size += s
				count += c
			}
		} else {
			size += uint64(f.Size())
			count++
		}
	}
	for _, f := range fi {
		if f.IsDir() && depth <= asyncDepth {
			next := <-nextOutput
			size += next.Size
			count += next.Count
		}
	}
	outputChan <- Output{Path: search, Size: size, Count: count}
	if depth <= maxDepth+1 {
		resultChan <- Output{Path: search, Size: size, Count: count}
	}
}
