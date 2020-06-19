package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/proto"
)

var out string

func init() {
	flag.StringVar(&out, "o", "/tmp/pprof.perfetto", "Path of converted perfetto file")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Usage:", os.Args[0], "PPROF_FILE")
		os.Exit(1)
	}

	// fmt.Println(flag.Arg(1))

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println("Cannot read pprof file:", err.Error())
		os.Exit(2)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	pprof, err := ParsePprof(data)
	if err != nil {
		panic(err)
	}

	profilePacket := convert(pprof)

	outData, err := proto.Marshal(profilePacket)
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(out)
	if err != nil {
		fmt.Println("cannot open output file:", err.Error())
		os.Exit(2)
	}
	defer outFile.Close()

	if _, err := outFile.Write(outData); err != nil {
		panic(err)
	}

	fmt.Println("done")
}
