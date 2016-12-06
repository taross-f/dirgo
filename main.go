package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const maxDepth int = 1

var asyncDepth = 2

type Output struct {
	Path  string
	Size  int64
	Count int64
}

var cpuCount int

func main() {
	cpuCount = runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)

	fmt.Println("start! Core:", cpuCount)
	root := os.Args[1]
	_, err := os.Stat(root)
	if err != nil {
		fmt.Println("You must set the valid target path.")
		panic(err)
	}

	// asyncDepth is optional
	if len(os.Args) >= 3 {
		ad := os.Args[2]
		argDepth, err := strconv.Atoi(ad)
		if err == nil {
			asyncDepth = argDepth
		}
	}
	fmt.Println("async depth: ", asyncDepth)

	checkNonRepeat(root)
}

func checkNonRepeat(root string) {
	syncStart := time.Now()
	paths := getTargetPaths(root, 0)
	fmt.Println("path count:" + fmt.Sprint(len(paths)))
	buf := ""

	result := make(chan Output, cpuCount)
	c := make(chan Output)
	go getSizeRecursiveNonRepeat(root, 0, c, result)
	for i := 0; i < len(paths); i++ {
		o := <-result
		buf += o.Path + "," + fmt.Sprint(o.Size) + "," + fmt.Sprint(o.Count) + "\n"
	}

	ioutil.WriteFile("./ouput_utf8.csv", []byte(buf), os.ModePerm)

	b, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(buf), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		fmt.Println(err.Error())
	}
	ioutil.WriteFile("./output.csv", b, os.ModePerm)
	syncEnd := time.Now()

	fmt.Println("---output: ", syncEnd.Sub(syncStart).Seconds(), "sec")
}

func getTargetPaths(root string, depth int) []string {
	fi, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Println("error occured: ", err.Error())
		return make([]string, 0) // if permission denied, return empty
		// panic(err)
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

func getSizeRecursive(root, search string) (int64, int64) {
	fi, err := ioutil.ReadDir(search)
	if err != nil {
		// fmt.Println("error occured: ", err.Error())
		return 0, 0 // if permission denied, return zeros
		// panic(err)
	}

	var size, count int64
	for _, f := range fi {
		if f.IsDir() {
			n := f.Name()
			s, c := getSizeRecursive(root, search+"/"+n)
			size += s
			count += c
		} else {
			size += f.Size()
			count++
		}
	}
	return size, count
}

func getSizeRecursiveNonRepeat(search string, depth int, outputChan chan Output, resultChan chan Output) {
	fi, err := ioutil.ReadDir(search)
	if err != nil {
		// fmt.Println("error occured: ", err.Error())
		outputChan <- Output{Path: search, Size: 0, Count: 0}
		if depth <= maxDepth+1 {
			resultChan <- Output{Path: search, Size: 0, Count: 0}
		}
		return
		// panic(err)
	}

	var size, count int64
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
			size += f.Size()
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
		fmt.Println(search, ",", size, ",", count)
	}
}
