package gpxreader

import (
    "fmt"
    "io"
    "os"
    "time"

    "github.com/dsoprea/go-gpx"
    "github.com/dsoprea/go-logging"
)

var (
    ErrNoTimestamps = fmt.Errorf("no points had timestamps")
)

type GpxSummary struct {
    Start time.Time
    Stop  time.Time
    Count int
}

func Summary(f io.Reader) (gs *GpxSummary, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    gs = new(GpxSummary)
    noTime := false

    tpc := func(tp *gpxcommon.TrackPoint) (err error) {
        defer func() {
            if state := recover(); state != nil {
                err = log.Wrap(state.(error))
            }
        }()

        if tp.Time.IsZero() {
            noTime = true
            return nil
        }

        gs.Count++

        if gs.Start.After(tp.Time) || gs.Start.IsZero() {
            gs.Start = tp.Time
        }

        if gs.Stop.Before(tp.Time) || gs.Stop.IsZero() {
            gs.Stop = tp.Time
        }

        return nil
    }

    if err := EnumerateTrackPoints(f, tpc); err != nil {
        log.Panic(err)
    }

    if gs.Count == 0 && noTime == true {
        log.Panic(ErrNoTimestamps)
    }

    return gs, nil
}

func FileSummary(filepath string) (gs *GpxSummary, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    f, err := os.Open(filepath)
    log.PanicIf(err)

    defer f.Close()

    gs, err = Summary(f)
    log.PanicIf(err)

    return gs, nil
}
