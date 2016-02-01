export GOPATH=${PWD}

GPXPARSE_EXECUTABLE_FILENAME=gpxparse

GPXREADER_FQ_PACKAGE=gpxreader/gpxreader
GPXREADER_SOURCEFILES=src/${GPXREADER_FQ_PACKAGE}/*.go
GPXPARSE_SOURCEFILES=${GPXREADER_SOURCEFILES} src/gpxreader/${GPXPARSE_EXECUTABLE_FILENAME}/*.go

.PHONY: all clean

all: bin/${GPXPARSE_EXECUTABLE_FILENAME}

clean:
	rm -fr bin pkg 

bin/${GPXPARSE_EXECUTABLE_FILENAME}: ${GPXPARSE_SOURCEFILES}
	go get gpxreader/${GPXPARSE_EXECUTABLE_FILENAME}
	go install gpxreader/${GPXPARSE_EXECUTABLE_FILENAME}
