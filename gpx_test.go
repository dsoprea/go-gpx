package gpxreader

import (
    "testing"
    "bytes"

    "github.com/dsoprea/go-logging"
)

type GpxPointCollector struct {
    FileVisits int
    FileVisitBalance int
    TrackVisits int
    TrackVisitBalance int
    TrackSegmentVisits int
    TrackSegmentVisitBalance int
    TrackPointVisits int
    TrackPointVisitBalance int
}

func NewGpsPointCollector() *GpxPointCollector {
    return new(GpxPointCollector)
}

func (gpc *GpxPointCollector) GpxOpen(gpx *Gpx) error {
    gpc.FileVisits++
    gpc.FileVisitBalance++

    return nil
}

func (gpc *GpxPointCollector) GpxClose(gpx *Gpx) error {
    gpc.FileVisitBalance--

    return nil
}

func (gpc *GpxPointCollector) TrackOpen(track *Track) error {
    gpc.TrackVisits++
    gpc.TrackVisitBalance++

    return nil
}

func (gpc *GpxPointCollector) TrackClose(track *Track) error {
    gpc.TrackVisitBalance--

    return nil
}

func (gpc *GpxPointCollector) TrackSegmentOpen(trackSegment *TrackSegment) error {
    gpc.TrackSegmentVisits++
    gpc.TrackSegmentVisitBalance++

    return nil
}

func (gpc *GpxPointCollector) TrackSegmentClose(trackSegment *TrackSegment) error {
    gpc.TrackSegmentVisitBalance--

    return nil
}

func (gpc *GpxPointCollector) TrackPointOpen(trackPoint *TrackPoint) error {
    gpc.TrackPointVisits++
    gpc.TrackPointVisitBalance++

    return nil
}

func (gpc *GpxPointCollector) TrackPointClose(trackPoint *TrackPoint) error {
    gpc.TrackPointVisitBalance--

    return nil
}

func TestFullGpxRead(t *testing.T) {
    b := bytes.NewBufferString(TestGpxData)
    gpc := NewGpsPointCollector()
    gp := NewGpxParser(b, gpc)

    if err := gp.Parse(); err != nil {
        log.Panic(err)
    }

    if gpc.FileVisits == 0 {
        t.Errorf("No file visits.")
    } else if gpc.FileVisitBalance != 0 {
        t.Errorf("File visits not balanced.")
    }

    if gpc.TrackVisits == 0 {
        t.Errorf("No track visits.")
    } else if gpc.TrackVisitBalance != 0 {
        t.Errorf("Track visits not balanced.")
    }

    if gpc.TrackSegmentVisits == 0 {
        t.Errorf("No track-segment visits.")
    } else if gpc.TrackSegmentVisitBalance != 0 {
        t.Errorf("Track-segment visits not balanced.")
    }

    if gpc.TrackPointVisits == 0 {
        t.Errorf("No track-point visits.")
    } else if gpc.TrackPointVisitBalance != 0 {
        t.Errorf("Track-point visits not balanced.")
    }

    if gpc.TrackPointVisits != 205 {
        t.Errorf("Points not correctly read.")
    }
}
