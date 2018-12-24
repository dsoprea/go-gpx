package gpxreader

import (
    "io"

    "github.com/dsoprea/go-logging"
)

type TrackPointCallback func(tp *TrackPoint) error

type SimpleGpxTrackVisitor struct {
    tpc TrackPointCallback
}

func NewSimpleGpxTrackVisitor(tpc TrackPointCallback) *SimpleGpxTrackVisitor {
    return &SimpleGpxTrackVisitor{
        tpc: tpc,
    }
}

func (gtv *SimpleGpxTrackVisitor) TrackPointOpen(tp *TrackPoint) (err error) {
    return nil
}

func (gtv *SimpleGpxTrackVisitor) TrackPointClose(tp *TrackPoint) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if err := gtv.tpc(tp); err != nil {
        log.Panic(err)
    }

    return nil
}

func EnumerateTrackPoints(f io.Reader, tpc TrackPointCallback) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    sgtv := NewSimpleGpxTrackVisitor(tpc)
    gp := NewGpxParser(f, sgtv)

    if err := gp.Parse(); err != nil {
        log.Panic(err)
    }

    return nil
}

func ExtractTrackPoints(f io.Reader) (points []TrackPoint, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    points = make([]TrackPoint, 0)

    tpc := func(tp *TrackPoint) (err error) {
        points = append(points, *tp)

        return nil
    }

    if err := EnumerateTrackPoints(f, tpc); err != nil {
        log.Panic(err)
    }

    return points, nil
}
