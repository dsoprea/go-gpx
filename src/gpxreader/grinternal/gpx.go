// GPX parser/visitor logic

package grinternal

import (
    "os"
    "strings"

    "encoding/xml"
)

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

type GpxParser struct {
    f *os.File
    decoder *xml.Decoder
    ns *Stack
    v GpxVisitor

    currentGpx *Gpx
    currentTrack *Track
    currentTrackSegment *TrackSegment
    currentTrackPoint *TrackPoint
}

// Create parser.
func NewGpxParser(filepath *string, visitor GpxVisitor) *GpxParser {
    f, err := os.Open(*filepath)
    if err != nil {
        panic(err)
    }

    decoder := xml.NewDecoder(f)
    ns := NewStack()

    return &GpxParser {
            f: f,
            decoder: decoder,
            ns: ns,
            v: visitor,
    }
}

// Close resources.
func (gp *GpxParser) Close() {
    gp.f.Close()
}

// Run the parse with a minimal memory footprint.
func (gp *GpxParser) Parse() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    for {
        token, err := gp.decoder.Token()
        if err != nil {
            break
        }
  
        switch t := token.(type) {
        case xml.StartElement:
            elmt := xml.StartElement(t)
            name := elmt.Name.Local

            gp.ns.Push(name)

            var attributes map[string]string = make(map[string]string)
            for _, a := range t.Attr {
                attributes[a.Name.Local] = a.Value
            }

            err := gp.handleStart(&name, &attributes)
            if err != nil {
                panic(err)
            }

        case xml.EndElement:
            gp.ns.Pop()

            elmt := xml.EndElement(t)
            name := elmt.Name.Local

            err := gp.handleEnd(&name)
            if err != nil {
                panic(err)
            }

        case xml.CharData:
            bytes := xml.CharData(t)
            s := strings.TrimSpace(string([]byte(bytes)))

            err := gp.handleCharData(&s)
            if err != nil {
                panic(err)
            }

        case xml.Comment:
        case xml.ProcInst:
        case xml.Directive:
        }
    }

    return nil
}

// Handle start tags.
func (gp *GpxParser) handleStart(tagName *string, attrp *map[string]string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    switch *tagName {
    case "gpx":
        err := gp.handleGpxEnd(attrp)
        if err != nil {
            panic(err)
        }

        gp.v.GpxOpen(gp.currentGpx)
    case "trk":
        gp.currentTrack = &Track {}
        gp.v.TrackOpen(gp.currentTrack)
    case "trkseg":
        gp.currentTrackSegment = &TrackSegment {}
        gp.v.TrackSegmentOpen(gp.currentTrackSegment)
    case "trkpt":
        err := gp.handleTrackPointEnd(attrp)
        if err != nil {
            panic(err)
        }

        gp.v.TrackPointOpen(gp.currentTrackPoint)
    }

    return nil
}

// Handle the end of a "GPX" [root] node.
func (gp *GpxParser) handleGpxEnd(attrp *map[string]string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    attr := *attrp

    gp.currentGpx = &Gpx {
            Xmlns: attr["xmlns"],
            Xsi: attr["xsi"],
            Creator: attr["creator"],
            SchemaLocation: attr["schemaLocation"],
    }

    versionRaw, ok := attr["version"]
    if ok == true {
        gp.currentGpx.Version = parseFloat32(versionRaw)
    }

    timeRaw, ok := attr["time"]
    if ok == true {
        gp.currentGpx.Time = parseIso8601Time(timeRaw)
    }

    return nil
}

// Handle the end of a track-point node.
func (gp *GpxParser) handleTrackPointEnd(attrp *map[string]string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    attr := *attrp

    gp.currentTrackPoint = &TrackPoint {}

    gp.currentTrackPoint.LatitudeDecimal = parseFloat64(attr["lat"])
    gp.currentTrackPoint.LongitudeDecimal = parseFloat64(attr["lon"])

    return nil
}

// Handle end tags.
func (gp *GpxParser) handleEnd(tagName *string) (err error) {
    switch *tagName {
    case "gpx":
        gp.v.GpxClose(gp.currentGpx)
        gp.currentGpx = nil
    case "trk":
        gp.v.TrackClose(gp.currentTrack)
        gp.currentTrack = nil
    case "trkseg":
        gp.v.TrackSegmentClose(gp.currentTrackSegment)
        gp.currentTrackSegment = nil
    case "trkpt":
        gp.v.TrackPointClose(gp.currentTrackPoint)
        gp.currentTrackPoint = nil
    }

    return nil
}

// Handle a string found between tags.
func (gp *GpxParser) handleCharData(s *string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    current := gp.ns.PeekFromEnd(0)
    parent := gp.ns.PeekFromEnd(1)

    if *s != "" && current != nil && parent != nil {
        currentName := current.(string)
        parentName := parent.(string)

        if parentName == "trkpt" {
            err := gp.handleTrackPointCharData(&currentName, s)
            if err != nil {
                panic(err)
            }
        }
    }

    return nil
}

// Handle values for the child nodes of a trackpoint node.
func (gp *GpxParser) handleTrackPointCharData(tagName *string, s *string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    switch *tagName {
    case "ele":
        gp.currentTrackPoint.Elevation = parseFloat32(*s)
    case "course":
        gp.currentTrackPoint.Course = parseFloat32(*s)
    case "speed":
        gp.currentTrackPoint.Speed = parseFloat32(*s)
    case "hdop":
        gp.currentTrackPoint.Hdop = parseFloat32(*s)
    case "src":
        gp.currentTrackPoint.Src = *s
    case "sat":
        gp.currentTrackPoint.SatelliteCount = parseUint8(*s)
    case "time":
        gp.currentTrackPoint.Time = parseIso8601Time(*s)
    }

    return nil
}
