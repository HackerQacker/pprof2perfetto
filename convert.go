package main

import (
	"time"

	"github.com/omerye/pprof2perfetto/protos/perfetto"
	"github.com/omerye/pprof2perfetto/protos/pprof"
)

/*
 * TODOs:
 * Verify `Iid`s purpose
 * Implement profile trace with memory (heap) data
 * Should make InternedString class with cache..
 */

func convert(p *pprof.Profile) *perfetto.Trace {
	return &perfetto.Trace{
		Packet: []*perfetto.TracePacket{convertProfile(p)},
	}
}

func convertProfile(p *pprof.Profile) *perfetto.TracePacket {
	// ts := uint64(time.Now().Unix())
	ts := uint64(p.TimeNanos)
	internedDataProxy := NewInternedDataProxy(p)
	return &perfetto.TracePacket{
		Timestamp:    &ts,
		InternedData: internedDataProxy.Get(),
		// Data: &perfetto.TracePacket_ProfilePacket{
		// 	ProfilePacket: &perfetto.ProfilePacket{
		// 		// Strings:    convertStringTable(p.StringTable),
		// 		// Mappings:   convertMappings(p.Mapping),
		// 		// Frames:     convertLocations(p.Location),
		// 		// Callstacks: convertSamples(p.Sample),
		// 		// TODO: continued looks god for me for relations between threads
		// 		// Continued: continued,
		// 	},
		// },
	}
}
