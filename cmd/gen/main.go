package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bigwhite/functrace/pkg/generator"
)

var (
	wrote bool
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
}

func usage() {
	fmt.Println("gen [-w] xxx.go")
	flag.PrintDefaults()
}

func main() {
	fmt.Println(os.Args)
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}

	if !strings.Contains(file, ".go") {
		usage()
		return
	}

	newSrc, err := generator.Rewrite(file)
	if err != nil {
		panic(err)
	}

	if newSrc == nil {
		// add nothing to the source file. no change
		fmt.Printf("no trace added for %s\n", file)
		return
	}

	if !wrote {
		fmt.Println(string(newSrc))
		return
	}

	// write to the source file
	f, err := os.OpenFile(file, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("open %s error: %v\n", file, err)
		return
	}
	defer f.Close()

	f.Seek(0, os.SEEK_SET)
	f.Truncate(0)
	f.Write(newSrc)
	fmt.Printf("add trace for %s ok\n", file)
}
