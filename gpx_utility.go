package gpxreader

import (
    "time"
    "os"
    "io"

    "github.com/dsoprea/go-logging"
)


type GpxSummary struct {
    Start time.Time
    Stop time.Time
    Count int
}

func Summary(f io.Reader) (gs *GpxSummary, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    gs = new(GpxSummary)

    tpc := func(tp *TrackPoint) (err error) {
        defer func() {
            if state := recover(); state != nil {
                err = log.Wrap(state.(error))
            }
        }()

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
