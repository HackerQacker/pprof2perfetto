package main

import (
	"errors"

	"github.com/omerye/pprof2perfetto/protos/perfetto"
	"github.com/omerye/pprof2perfetto/protos/pprof"
)

func convertID(pprofID int64) *uint64 {
	perfettoID := uint64(pprofID)
	return &perfettoID
}

func convertIDs(pprofIDs []int64) []uint64 {
	perfettoIDs := make([]uint64, len(pprofIDs))
	for i, id := range pprofIDs {
		perfettoIDs[i] = uint64(id)
	}
	return perfettoIDs
}

type InternedStringProxy struct {
	internedStrings []*perfetto.InternedString
	stringsCache    map[string]uint64
	stringTable     []string
}

func NewInternedStringProxy(stringTable []string) *InternedStringProxy {
	return &InternedStringProxy{
		stringsCache: make(map[string]uint64),
		stringTable:  stringTable,
	}
}

func (this *InternedStringProxy) FromST(i int64) uint64 {
	if i < 0 || int(i) >= len(this.stringTable) {
		panic("invalid string table id")
	}

	str := this.stringTable[i]
	return this.fromString(str)
}

func (this *InternedStringProxy) fromString(str string) uint64 {
	if fromCache, exists := this.stringsCache[str]; exists {
		return fromCache
	}

	iid := uint64(len(this.internedStrings))
	this.internedStrings = append(this.internedStrings, &perfetto.InternedString{
		Iid: &iid,
		Str: []byte(str),
	})
	return iid
}

func (this *InternedStringProxy) PFromST(i int64) *uint64 {
	iid := this.FromST(i)
	return &iid
}

func (this *InternedStringProxy) Get() []*perfetto.InternedString {
	return this.internedStrings
}

type InternedDataProxy struct {
	p            *pprof.Profile
	internedData *perfetto.InternedData
}

func NewInternedDataProxy(p *pprof.Profile) *InternedDataProxy {
	return &InternedDataProxy{
		p:            p,
		internedData: makeInternedData(p),
	}
}

func makeInternedData(p *pprof.Profile) *perfetto.InternedData {
	/* TODOs:
	 * Parse comments to DebugAnnotations
	 * Parse labels to categories
	 */

	buildIds := NewInternedStringProxy(p.StringTable)
	mappingPaths := NewInternedStringProxy(p.StringTable)
	sourcePaths := NewInternedStringProxy(p.StringTable)
	functionNames := NewInternedStringProxy(p.StringTable)

	// categories := make([]*perfetto.EventCategory, len(p.SampleType))
	// for i, v := range p.SampleType {
	// 	iid := uint64(i)
	// 	name := p.StringTable[v.Type]
	// 	categories[i] = &perfetto.EventCategory{
	// 		Iid:  &iid,
	// 		Name: &name,
	// 	}
	// }

	mappings := make([]*perfetto.Mapping, len(p.Mapping))
	for i, m := range p.Mapping {
		iid := uint64(i)
		mappings[i] = &perfetto.Mapping{
			Iid:           &iid,
			BuildId:       buildIds.PFromST(m.BuildId),
			ExactOffset:   &m.FileOffset,
			StartOffset:   &m.FileOffset,
			Start:         &m.MemoryStart,
			End:           &m.MemoryLimit,
			PathStringIds: []uint64{mappingPaths.FromST(m.Filename)},
		}
	}

	frames := make([]*perfetto.Frame, len(p.Location))
	var sourceLocations []*perfetto.SourceLocation
	var funcNameID, fileNameID int64
	for i, l := range p.Location {
		for _, line := range l.Line {
			iid := uint64(len(sourceLocations))
			function, err := findFunction(p, line.FunctionId)
			if err != nil {
				panic(err)
			}

			funcName := p.StringTable[function.Name]
			fileName := p.StringTable[function.Filename]
			lineNum := uint32(line.Line)
			sourceLocations = append(sourceLocations, &perfetto.SourceLocation{
				Iid:          &iid,
				FileName:     &fileName,
				FunctionName: &funcName,
				LineNumber:   &lineNum,
			})
			funcNameID = function.Name
			fileNameID = function.Filename
		}

		sourcePaths.PFromST(fileNameID) // Just to set it in.
		iid := uint64(i)
		frames[i] = &perfetto.Frame{
			Iid:            &iid,
			FunctionNameId: functionNames.PFromST(funcNameID),
			MappingId:      &l.MappingId, // It should be the same id (index)
			RelPc:          &l.Address,
		}
	}

	callstacks := make([]*perfetto.Callstack, len(p.Sample))
	for i, c := range p.Sample {
		iid := uint64(i)
		callstacks[i] = &perfetto.Callstack{
			Iid:      &iid,
			FrameIds: c.LocationId,
		}
	}

	return &perfetto.InternedData{
		Callstacks:      callstacks,
		BuildIds:        buildIds.Get(),
		Mappings:        mappings,
		MappingPaths:    mappingPaths.Get(),
		Frames:          frames,
		FunctionNames:   functionNames.Get(),
		SourcePaths:     sourcePaths.Get(),
		SourceLocations: sourceLocations,
	}
}

func findFunction(p *pprof.Profile, id uint64) (*pprof.Function, error) {
	for _, f := range p.Function {
		if f.Id == id {
			return f, nil
		}
	}
	return nil, errors.New("invalid function id")
}

func findLocation(p *pprof.Profile, id uint64) (*pprof.Location, error) {
	for _, l := range p.Location {
		if l.Id == id {
			return l, nil
		}
	}

	return nil, errors.New("invalid location id")
}

func (this *InternedDataProxy) Get() *perfetto.InternedData {
	return this.internedData
}

func (this *InternedDataProxy) Callstack(i uint64) (*perfetto.Callstack, error) {
	if int(i) >= len(this.internedData.Callstacks) {
		return nil, errors.New("callstack for found")
	}

	return this.internedData.Callstacks[i], nil
}
