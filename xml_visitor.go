package gpxreader

import (
    "time"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-xmlvisitor/xmlvisitor"
)


type GpxFileVisitor interface {
    GpxOpen(g *Gpx) error
    GpxClose(g *Gpx) error
}


type GpxTrackVisitor interface {
    TrackOpen(t *Track) error
    TrackClose(t *Track) error
}


type GpxTrackSegmentVisitor interface {
    TrackSegmentOpen(ts *TrackSegment) error
    TrackSegmentClose(ts *TrackSegment) error
}


type GpxTrackPointVisitor interface {
    TrackPointOpen(tp *TrackPoint) error
    TrackPointClose(tp *TrackPoint) error
}


type xmlVisitor struct {
    gp *GpxParser
    v interface{}

    currentGpx *Gpx
    currentTrack *Track
    currentTrackSegment *TrackSegment
    currentTrackPoint *TrackPoint
}

func newXmlVisitor(gp *GpxParser, v interface{}) (*xmlVisitor) {
    return &xmlVisitor {
        gp: gp,
        v: v,
    }
}

func (xv *xmlVisitor) HandleStart(tagName *string, attrp *map[string]string, xp *xmlvisitor.XmlParser) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    switch *tagName {
    case "gpx":
        if err := xv.handleGpxStart(attrp); err != nil {
            log.Panic(err)
        }

        if gfv, ok := xv.v.(GpxFileVisitor); ok == true {
            if err := gfv.GpxOpen(xv.currentGpx); err != nil {
                log.Panic(err)
            }
        }
    case "trk":
        xv.currentTrack = new(Track)

        if gtv, ok := xv.v.(GpxTrackVisitor); ok == true {
            if err := gtv.TrackOpen(xv.currentTrack); err != nil {
                log.Panic(err)
            }
        }
    case "trkseg":
        xv.currentTrackSegment = new(TrackSegment)

        if gtsv, ok := xv.v.(GpxTrackSegmentVisitor); ok == true {
            if err := gtsv.TrackSegmentOpen(xv.currentTrackSegment); err != nil {
                log.Panic(err)
            }
        }
    case "trkpt":
        if err := xv.handleTrackPointEnd(attrp); err != nil {
            log.Panic(err)
        }

        if gtpv, ok := xv.v.(GpxTrackPointVisitor); ok == true {
            if err := gtpv.TrackPointOpen(xv.currentTrackPoint); err != nil {
                log.Panic(err)
            }
        }
    }

    return nil
}

func (xv *xmlVisitor) HandleEnd(tagName *string, xp *xmlvisitor.XmlParser) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    switch *tagName {
    case "gpx":
        if gfv, ok := xv.v.(GpxFileVisitor); ok == true {
            if err := gfv.GpxClose(xv.currentGpx); err != nil {
                log.Panic(err)
            }
        }

        xv.currentGpx = nil
    case "trk":
        if gtv, ok := xv.v.(GpxTrackVisitor); ok == true {
            if err := gtv.TrackClose(xv.currentTrack); err != nil {
                log.Panic(err)
            }
        }

        xv.currentTrack = nil
    case "trkseg":
        if gtsv, ok := xv.v.(GpxTrackSegmentVisitor); ok == true {
            if err := gtsv.TrackSegmentClose(xv.currentTrackSegment); err != nil {
                log.Panic(err)
            }
        }

        xv.currentTrackSegment = nil
    case "trkpt":
        if gtpv, ok := xv.v.(GpxTrackPointVisitor); ok == true {
            if err := gtpv.TrackPointClose(xv.currentTrackPoint); err != nil {
                log.Panic(err)
            }
        }

        xv.currentTrackPoint = nil
    }

    return nil
}

func (xv *xmlVisitor) HandleValue(tagName *string, value *string, xp *xmlvisitor.XmlParser) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    ns := xp.NodeStack()
    parent := ns.PeekFromEnd(0)

    if parent != nil {
        parentName := parent.(string)

        if parentName == "trkpt" {
            if err := xv.handleTrackPointValue(tagName, value); err != nil {
                log.Panic(err)
            }
        }
    }

    return nil
}

// Parse the 8601 timestamps.
func (xv *xmlVisitor) parseTimestamp(phrase *string) (timestamp time.Time, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()
    
    t, err := time.Parse(time.RFC3339Nano, *phrase)
    log.PanicIf(err)

    return t, nil
}

// Handle the end of a "GPX" [root] node.
func (xv *xmlVisitor) handleGpxStart(attrp *map[string]string) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    attr := *attrp

    xv.currentGpx = &Gpx {
            Xmlns: attr["xmlns"],
            Xsi: attr["xsi"],
            Creator: attr["creator"],
            SchemaLocation: attr["schemaLocation"],
    }

    versionRaw, ok := attr["version"]
    if ok == true {
        xv.currentGpx.Version = parseFloat32(versionRaw)
    }

    timeRaw, ok := attr["time"]
    if ok == true {
        xv.currentGpx.Time, err = xv.parseTimestamp(&timeRaw)
        log.PanicIf(err)
    }

    return nil
}

// Handle the end of a track-point node.
func (xv *xmlVisitor) handleTrackPointEnd(attrp *map[string]string) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    attr := *attrp

    xv.currentTrackPoint = &TrackPoint {}

    xv.currentTrackPoint.LatitudeDecimal = parseFloat64(attr["lat"])
    xv.currentTrackPoint.LongitudeDecimal = parseFloat64(attr["lon"])

    return nil
}

// Handle values for the child nodes of a trackpoint node.
func (xv *xmlVisitor) handleTrackPointValue(tagName *string, s *string) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    switch *tagName {
    case "ele":
        xv.currentTrackPoint.Elevation = parseFloat32(*s)
    case "course":
        xv.currentTrackPoint.Course = parseFloat32(*s)
    case "speed":
        xv.currentTrackPoint.Speed = parseFloat32(*s)
    case "hdop":
        xv.currentTrackPoint.Hdop = parseFloat32(*s)
    case "src":
        xv.currentTrackPoint.Src = *s
    case "sat":
        xv.currentTrackPoint.SatelliteCount = parseUint8(*s)
    case "time":
        xv.currentTrackPoint.Time, err = xv.parseTimestamp(s)
        log.PanicIf(err)
    }

    return nil
}
