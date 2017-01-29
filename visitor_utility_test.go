package gpxreader

import (
    "testing"
    "bytes"

    "github.com/dsoprea/go-logging"
)

func TestEnumerateTrackPoints(t *testing.T) {
    n := 0
    cb := func(tp *TrackPoint) error {
        n++

        return nil
    }

    b := bytes.NewBufferString(TestGpxData)

    if err := EnumerateTrackPoints(b, cb); err != nil {
        log.Panic(err)
    }

    if n != 205 {
        t.Errorf("Points not read correctly.")
    }
}

func TestExtractTrackPoints(t *testing.T) {
    b := bytes.NewBufferString(TestGpxData)
    points, err := ExtractTrackPoints(b)
    log.PanicIf(err)

    if len(points) != 205 {
        t.Errorf("Points not read correctly.")
    }
}
