package gpxwriter

import (
    "bytes"
    "fmt"
    "testing"
    "time"

    "github.com/dsoprea/go-logging"
)

func TestBuilder_Gpx(t *testing.T) {
    buffer := new(bytes.Buffer)

    b := NewBuilder(buffer)
    gb := b.Gpx()
    gb.EndGpx()

    expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd"></gpx>`

    if buffer.String() != expected {
        fmt.Printf("\nACTUAL:\n%s\n", buffer.String())
        fmt.Printf("\nEXPECTED:\n%s\n", expected)

        t.Fatalf("Output not expected.")
    }
}

func TestBuilder_Track(t *testing.T) {
    buffer := new(bytes.Buffer)

    b := NewBuilder(buffer)
    gb := b.Gpx()

    tb, err := gb.Track()
    log.PanicIf(err)

    err = tb.EndTrack()
    log.PanicIf(err)

    gb.EndGpx()

    expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
  <trk></trk>
</gpx>`

    if buffer.String() != expected {
        t.Fatalf("Output not expected:\n%s", buffer.String())
    }
}

func TestBuilder_TrackSegment(t *testing.T) {
    buffer := new(bytes.Buffer)

    b := NewBuilder(buffer)
    gb := b.Gpx()

    tb, err := gb.Track()
    log.PanicIf(err)

    tsb, err := tb.TrackSegment()
    log.PanicIf(err)

    err = tsb.EndTrackSegment()
    log.PanicIf(err)

    err = tb.EndTrack()
    log.PanicIf(err)

    gb.EndGpx()

    expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
  <trk>
    <trkseg></trkseg>
  </trk>
</gpx>`

    if buffer.String() != expected {
        t.Fatalf("Output not expected:\n%s", buffer.String())
    }
}

func TestBuilder_TrackPoint(t *testing.T) {
    buffer := new(bytes.Buffer)

    b := NewBuilder(buffer)
    gb := b.Gpx()

    tb, err := gb.Track()
    log.PanicIf(err)

    tsb, err := tb.TrackSegment()
    log.PanicIf(err)

    tpb := tsb.TrackPoint()

    tpb.LatitudeDecimal = .123
    tpb.LongitudeDecimal = .456

    now := time.Now()
    tpb.Time = now

    err = tpb.Write()
    log.PanicIf(err)

    err = tsb.EndTrackSegment()
    log.PanicIf(err)

    err = tb.EndTrack()
    log.PanicIf(err)

    gb.EndGpx()

    expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
  <trk>
    <trkseg>
      <trkpt lat="0.123" lon="0.456">
        <time>` + now.UTC().Format("2006-01-02T15:04:05-0700") + `</time>
      </trkpt>
    </trkseg>
  </trk>
</gpx>`

    if buffer.String() != expected {
        fmt.Printf("\nACTUAL:\n%s\n", buffer.String())
        fmt.Printf("\nEXPECTED:\n%s\n", expected)
    }
}

func TestBuilder_TrackPoint_Multiple(t *testing.T) {
    buffer := new(bytes.Buffer)

    b := NewBuilder(buffer)
    gb := b.Gpx()

    tb, err := gb.Track()
    log.PanicIf(err)

    tsb, err := tb.TrackSegment()
    log.PanicIf(err)

    // Point 1

    tpb := tsb.TrackPoint()

    tpb.LatitudeDecimal = .123
    tpb.LongitudeDecimal = .456

    now1 := time.Now()
    tpb.Time = now1

    err = tpb.Write()
    log.PanicIf(err)

    // Point 2

    tpb = tsb.TrackPoint()

    tpb.LatitudeDecimal = .123
    tpb.LongitudeDecimal = .456

    now2 := time.Now()
    tpb.Time = now2

    err = tpb.Write()
    log.PanicIf(err)

    // Point 3

    tpb = tsb.TrackPoint()

    tpb.LatitudeDecimal = .123
    tpb.LongitudeDecimal = .456

    now3 := time.Now()
    tpb.Time = now3

    err = tpb.Write()
    log.PanicIf(err)

    err = tsb.EndTrackSegment()
    log.PanicIf(err)

    err = tb.EndTrack()
    log.PanicIf(err)

    gb.EndGpx()

    expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
  <trk>
    <trkseg>
      <trkpt lat="0.123" lon="0.456">
        <time>` + now1.UTC().Format("2006-01-02T15:04:05-0700") + `</time>
      </trkpt>
      <trkpt lat="0.123" lon="0.456">
        <time>` + now2.UTC().Format("2006-01-02T15:04:05-0700") + `</time>
      </trkpt>
      <trkpt lat="0.123" lon="0.456">
        <time>` + now3.UTC().Format("2006-01-02T15:04:05-0700") + `</time>
      </trkpt>
    </trkseg>
  </trk>
</gpx>`

    if buffer.String() != expected {
        fmt.Printf("\nACTUAL:\n%s\n", buffer.String())
        fmt.Printf("\nEXPECTED:\n%s\n", expected)
    }
}
