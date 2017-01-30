package gpxreader

import (
    "testing"
    "bytes"
    "time"

    "github.com/dsoprea/go-logging"
)

func TestSummary(t *testing.T) {
    b := bytes.NewBufferString(TestGpxData)

    gs, err := Summary(b)
    log.PanicIf(err)

    if gs.Start.Format(time.RFC3339) != "2016-12-02T08:05:44Z" {
        t.Errorf("Start time is not correct.")
    } else if gs.Stop.Format(time.RFC3339) != "2016-12-03T07:57:07Z" {
        t.Errorf("Stop time is not correct.")
    } else if gs.Count != 205 {
        t.Errorf("Point count is not correct.")
    }
}
