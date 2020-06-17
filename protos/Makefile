
REMOTE_PPROF_PROTO=https://github.com/google/pprof/raw/master/proto/profile.proto
REMOTE_PERFETTO_PROTO=https://github.com/google/perfetto/raw/master/protos/perfetto/trace/perfetto_trace.proto

PROTOS = pprof/profile.proto perfetto/perfetto_trace.proto
# PROTOS = $(shell find . -type f -name *.proto)
SRC = $(PROTOS:%.proto=%.pb.go)

# all: pprof/profile.pb.go
all: $(SRC)

%.pb.go: PKG_NAME=$(shell basename $(dir $@))
%.pb.go: %.proto
	echo 'option go_package = "$(PKG_NAME)";' >> $<
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative $<

perfetto/perfetto_trace.proto:
	wget ${REMOTE_PERFETTO_PROTO} -O $@


pprof/profile.proto:
	wget ${REMOTE_PPROF_PROTO} -O $@

clean:
	$(RM) $(SRC)

clean-all: clean
clean-all:
	$(RM) $(PROTOS)
