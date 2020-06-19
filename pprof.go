package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/omerye/pprof2perfetto/protos/pprof"
	"google.golang.org/protobuf/proto"
)

func ParsePprof(data []byte) (*pprof.Profile, error) {
	if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		gz, err := gzip.NewReader(bytes.NewBuffer(data))
		if err == nil {
			data, err = ioutil.ReadAll(gz)
		}
		if err != nil {
			return nil, fmt.Errorf("decompressing profile: %v", err)
		}
	}

	return ParseUncompressedPprof(data)
}

func ParseUncompressedPprof(data []byte) (*pprof.Profile, error) {
	p := new(pprof.Profile)
	err := proto.Unmarshal(data, p)
	return p, err
}
