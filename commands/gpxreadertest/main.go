package main

import (
    "os"
    "fmt"

    flags "github.com/jessevdk/go-flags"

    "github.com/dsoprea/go-gpxreader/gpxreader"
)

type gpxVisitor struct {}

func newGpxVisitor() *gpxVisitor {
    return &gpxVisitor {}
}

func (gv *gpxVisitor) GpxOpen(gpx *gpxreader.Gpx) error {
    fmt.Printf("GPX: %s\n", gpx)

    return nil
}

func (gv *gpxVisitor) GpxClose(gpx *gpxreader.Gpx) error {
    return nil
}

func (gv *gpxVisitor) TrackOpen(track *gpxreader.Track) error {
    fmt.Printf("Track: %s\n", track)

    return nil
}

func (gv *gpxVisitor) TrackClose(track *gpxreader.Track) error {
    return nil
}

func (gv *gpxVisitor) TrackSegmentOpen(trackSegment *gpxreader.TrackSegment) error {
    fmt.Printf("Track segment: %s\n", trackSegment)

    return nil
}

func (gv *gpxVisitor) TrackSegmentClose(trackSegment *gpxreader.TrackSegment) error {
    return nil
}

func (gv *gpxVisitor) TrackPointOpen(trackPoint *gpxreader.TrackPoint) error {
    return nil
}

func (gv *gpxVisitor) TrackPointClose(trackPoint *gpxreader.TrackPoint) error {
    fmt.Printf("Point: %s\n", trackPoint)

    return nil
}

type options struct {
    GpxFilepath string  `short:"f" long:"gpx-filepath" description:"GPX file-path" required:"true"`
}

func readOptions () *options {
    o := options {}

    _, err := flags.Parse(&o)
    if err != nil {
        os.Exit(1)
    }

    return &o
}

func main() {
    var gpxFilepath string

    o := readOptions()

    gpxFilepath = o.GpxFilepath

    f, err := os.Open(gpxFilepath)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    gv := newGpxVisitor()
    gp := gpxreader.NewGpxParser(f, gv)

    err = gp.Parse()
    if err != nil {
        print("Error: %s\n", err.Error())
        os.Exit(1)
    }
}
