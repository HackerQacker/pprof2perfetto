package main

import (
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
	trace := new(perfetto.Trace)
	addStreamProfile(p, trace)
	return trace
}

// func convertProfile(p *pprof.Profile) *perfetto.TracePacket {
// 	// ts := uint64(time.Now().Unix())
// 	ts := uint64(p.TimeNanos)
// 	internedDataProxy := NewInternedDataProxy(p)
// 	return &perfetto.TracePacket{
// 		Timestamp:    &ts,
// 		InternedData: internedDataProxy.Get(),
// 		Data:         makeTrackEventData(p, internedDataProxy),
// 	}
// }

func addStreamProfile(p *pprof.Profile, trace *perfetto.Trace) {
	ts := uint64(p.TimeNanos)
	internedDataProxy := NewInternedDataProxy(p)
	callstacks := make([]uint64, len(p.Sample))
	for i := range p.Sample {
		callstacks[i] = uint64(i)
	}

	packet := &perfetto.TracePacket{
		Timestamp:    &ts,
		InternedData: internedDataProxy.Get(),
		Data: &perfetto.TracePacket_StreamingProfilePacket{
			StreamingProfilePacket: &perfetto.StreamingProfilePacket{
				CallstackIid: callstacks,
				// TimestampDeltaUs: &ts, // ??
			},
		},
	}

	trace.Packet = append(trace.Packet, packet)
}

// func makeProfilePacket(p *pprof.Profile, interned *InternedDataProxy) *perfetto.TracePacket_ProfilePacket {
// 	return &perfetto.TracePacket_ProfilePacket{
// 		&perfetto.ProfilePacket{
// 			ProcessDumps: []*perfetto.ProfilePacket_ProcessHeapSamples{
// 				&perfetto.ProfilePacket_ProcessHeapSamples{
// 					Pid:       pid,
// 					Samples:   []*perfetto.ProfilePacket_HeapSample{},
// 					Timestamp: timestamp,
// 					Stats:     &perfetto.ProfilePacket_ProcessStats{},
// 				},
// 			},
// 		}}
// }

// func convertToTrackTrace(p *pprof.Profile) *perfetto.Trace {
// 	ts := uint64(p.TimeNanos)
// 	internedDataProxy := NewInternedDataProxy(p)

// 	var packets []*perfetto.TracePacket
// 	// packets := make([]*perfetto.TracePacket, len(p.Sample))
// 	for _, sample := range p.Sample {
// 		packet := &perfetto.TracePacket{
// 		Timestamp: &ts,
// 		InternedData: internedDataProxy.Get(),
// 		Data: &,
// 	}
// 	}
// }
// func makeTrackEventData(p *pprof.Profile, internedDataProxy *InternedDataProxy) *perfetto.TracePacket_TrackEvent {
// 	return &perfetto.TracePacket_TrackEvent{
// 		TrackEvent: &perfetto.TrackEvent{

// 		},
// 	}
// }
