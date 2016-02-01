package main

import (
    "os"
    "fmt"

    "gpxreader/grinternal"
)

type GpxVisitor struct {}

func NewGpxVisitor() (*GpxVisitor) {
    return &GpxVisitor {}
}

func (gv *GpxVisitor) GpxOpen(gpx *grinternal.Gpx) error {
    fmt.Printf("GPX: %s\n", *gpx)

    return nil
}

func (gv *GpxVisitor) GpxClose(gpx *grinternal.Gpx) error {
    return nil
}

func (gv *GpxVisitor) TrackOpen(track *grinternal.Track) error {
    return nil
}

func (gv *GpxVisitor) TrackClose(track *grinternal.Track) error {
    return nil
}

func (gv *GpxVisitor) TrackSegmentOpen(trackSegment *grinternal.TrackSegment) error {
    return nil
}

func (gv *GpxVisitor) TrackSegmentClose(trackSegment *grinternal.TrackSegment) error {
    return nil
}

func (gv *GpxVisitor) TrackPointOpen(trackPoint *grinternal.TrackPoint) error {
    return nil
}

func (gv *GpxVisitor) TrackPointClose(trackPoint *grinternal.TrackPoint) error {
    fmt.Printf("Point: %s\n", trackPoint)

    return nil
}

func main() {
    filepath := "test/20130729.gpx"
    gv := NewGpxVisitor()
    gp := grinternal.NewGpxParser(&filepath, gv)

    defer gp.Close()

    err := gp.Parse()
    if err != nil {
        print("Error: %s\n", err.Error())
        os.Exit(1)
    }
}
