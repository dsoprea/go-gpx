// GPX parser/visitor logic

package gpxreader

import (

// TODO(dustin): Implement the URL.
    "xmlvisitor/xmlvisitor"

)

type GpxParser struct {
    xp *xmlvisitor.XmlParser
}

// Create parser.
func NewGpxParser(filepath *string, visitor GpxVisitor) *GpxParser {
    gp := &GpxParser {}

    v := newXmlVisitor(gp, visitor)
    gp.xp = xmlvisitor.NewXmlParser(filepath, v)

    return gp
}

// Close resources.
func (gp *GpxParser) Close() {
    gp.xp.Close()
}

// Run the parse with a minimal memory footprint.
func (gp *GpxParser) Parse() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    err = gp.xp.Parse()
    if err != nil {
        panic(err)
    }

    return nil
}

type GpxVisitor interface {
    GpxOpen(gpx *Gpx) error
    GpxClose(gpx *Gpx) error
    TrackOpen(track *Track) error
    TrackClose(track *Track) error
    TrackSegmentOpen(trackSegment *TrackSegment) error
    TrackSegmentClose(trackSegment *TrackSegment) error
    TrackPointOpen(trackPoint *TrackPoint) error
    TrackPointClose(trackPoint *TrackPoint) error
}

type xmlVisitor struct {
    gp *GpxParser
    v GpxVisitor

    currentGpx *Gpx
    currentTrack *Track
    currentTrackSegment *TrackSegment
    currentTrackPoint *TrackPoint
}

func newXmlVisitor(gp *GpxParser, v GpxVisitor) (*xmlVisitor) {
    return &xmlVisitor {
            gp: gp,
            v: v,
    }
}

func (xv *xmlVisitor) HandleStart(tagName *string, attrp *map[string]string, xp *xmlvisitor.XmlParser) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    switch *tagName {
    case "gpx":
        err := xv.handleGpxStart(attrp)
        if err != nil {
            panic(err)
        }

        err = xv.v.GpxOpen(xv.currentGpx)
        if err != nil {
            panic(err)
        }
    case "trk":
        xv.currentTrack = &Track {}

        err := xv.v.TrackOpen(xv.currentTrack)
        if err != nil {
            panic(err)
        }
    case "trkseg":
        xv.currentTrackSegment = &TrackSegment {}

        err := xv.v.TrackSegmentOpen(xv.currentTrackSegment)
        if err != nil {
            panic(err)
        }
    case "trkpt":
        err := xv.handleTrackPointEnd(attrp)
        if err != nil {
            panic(err)
        }

        err = xv.v.TrackPointOpen(xv.currentTrackPoint)
        if err != nil {
            panic(err)
        }
    }

    return nil
}

func (xv *xmlVisitor) HandleEnd(tagName *string, xp *xmlvisitor.XmlParser) error {
    switch *tagName {
    case "gpx":
        
        err := xv.v.GpxClose(xv.currentGpx)
        if err != nil {
            panic(err)
        }

        xv.currentGpx = nil

    case "trk":
        
        err := xv.v.TrackClose(xv.currentTrack)
        if err != nil {
            panic(err)
        }

        xv.currentTrack = nil

    case "trkseg":
        
        err := xv.v.TrackSegmentClose(xv.currentTrackSegment)
        if err != nil {
            panic(err)
        }

        xv.currentTrackSegment = nil

    case "trkpt":
        
        err := xv.v.TrackPointClose(xv.currentTrackPoint)
        if err != nil {
            panic(err)
        }

        xv.currentTrackPoint = nil
    }

    return nil
}

func (xv *xmlVisitor) HandleCharData(data *string, xp *xmlvisitor.XmlParser) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

//    if xp.LastState() != xmlvisitor.XmlPartStartTag {
//        return nil
//    }

    ns := xp.NodeStack()
    current := ns.PeekFromEnd(0)
    parent := ns.PeekFromEnd(1)

    if *data != "" && current != nil && parent != nil {
        currentName := current.(string)
        parentName := parent.(string)

        if parentName == "trkpt" {
            err := xv.handleTrackPointCharData(&currentName, data)
            if err != nil {
                panic(err)
            }
        }
    }

    return nil
}


// Handle the end of a "GPX" [root] node.
func (xv *xmlVisitor) handleGpxStart(attrp *map[string]string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
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
        xv.currentGpx.Time = parseIso8601Time(timeRaw)
    }

    return nil
}

// Handle the end of a track-point node.
func (xv *xmlVisitor) handleTrackPointEnd(attrp *map[string]string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    attr := *attrp

    xv.currentTrackPoint = &TrackPoint {}

    xv.currentTrackPoint.LatitudeDecimal = parseFloat64(attr["lat"])
    xv.currentTrackPoint.LongitudeDecimal = parseFloat64(attr["lon"])

    return nil
}

// Handle values for the child nodes of a trackpoint node.
func (xv *xmlVisitor) handleTrackPointCharData(tagName *string, s *string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
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
        xv.currentTrackPoint.Time = parseIso8601Time(*s)
    }

    return nil
}
