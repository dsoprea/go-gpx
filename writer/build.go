package gpxwriter

import (
    "io"
    "strconv"
    "time"

    "encoding/xml"

    "github.com/dsoprea/go-logging"
)

type Builder struct {
    w       io.Writer
    encoder *xml.Encoder
}

func NewBuilder(w io.Writer) *Builder {
    w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))

    encoder := xml.NewEncoder(w)

    encoder.Indent("", "  ")

    return &Builder{
        w:       w,
        encoder: encoder,
    }
}

type GpxBuilder struct {
    b *Builder
}

func (b *Builder) Gpx() *GpxBuilder {

    // Add <gpx> tag:
    //
    // <gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" creator="Oregon 400t" version="1.1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">//     // <gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" creator="Oregon 400t" version="1.1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">

    attrs := []xml.Attr{
        {
            Name:  xml.Name{"", "xmlns"},
            Value: "http://www.topografix.com/GPX/1/1",
        },
        {
            Name:  xml.Name{"", "xmlns:gpxx"},
            Value: "http://www.garmin.com/xmlschemas/GpxExtensions/v3",
        },
        {
            Name:  xml.Name{"", "gpxtpx"},
            Value: "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
        },
        {
            Name:  xml.Name{"", "xmlns:xsi"},
            Value: "http://www.w3.org/2001/XMLSchema-instance",
        },
        {
            Name:  xml.Name{"", "xsi:schemaLocation"},
            Value: "http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd",
        },
    }

    gpxStart := xml.StartElement{
        Name: xml.Name{
            Space: "",
            Local: "gpx",
        },
        Attr: attrs,
    }

    err := b.encoder.EncodeToken(gpxStart)
    log.PanicIf(err)

    return &GpxBuilder{
        b: b,
    }
}

func (gb *GpxBuilder) EndGpx() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    endElement := xml.EndElement{
        Name: xml.Name{
            Space: "",
            Local: "gpx",
        },
    }

    err = gb.b.encoder.EncodeToken(endElement)
    log.PanicIf(err)

    gb.b.encoder.Flush()

    return nil
}

type GpxTrackBuilder struct {
    b *Builder
}

func (gb *GpxBuilder) Track() (gpb *GpxTrackBuilder, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    // Add <gpx> tag:
    //
    // <trk>

    trkStart := xml.StartElement{
        Name: xml.Name{
            Space: "",
            Local: "trk",
        },
    }

    err = gb.b.encoder.EncodeToken(trkStart)
    log.PanicIf(err)

    gtb := &GpxTrackBuilder{
        b: gb.b,
    }

    return gtb, nil
}

func (gtb *GpxTrackBuilder) EndTrack() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    endElement := xml.EndElement{
        Name: xml.Name{
            Space: "",
            Local: "trk",
        },
    }

    err = gtb.b.encoder.EncodeToken(endElement)
    log.PanicIf(err)

    return nil
}

type GpxTrackSegmentBuilder struct {
    b *Builder
}

func (gtb *GpxTrackBuilder) TrackSegment() (gtsb *GpxTrackSegmentBuilder, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    // Add <trkseg> tag:
    //
    // <trkseg>

    trksegStart := xml.StartElement{
        Name: xml.Name{
            Space: "",
            Local: "trkseg",
        },
    }

    err = gtb.b.encoder.EncodeToken(trksegStart)
    log.PanicIf(err)

    gtsb = &GpxTrackSegmentBuilder{
        b: gtb.b,
    }

    return gtsb, nil
}

func (gtsb *GpxTrackSegmentBuilder) EndTrackSegment() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    endElement := xml.EndElement{
        Name: xml.Name{
            Space: "",
            Local: "trkseg",
        },
    }

    err = gtsb.b.encoder.EncodeToken(endElement)
    log.PanicIf(err)

    return nil
}

type GpxTrackPointBuilder struct {
    b *Builder

    // TODO(dustin): !! Should we just marshall this type directly?
    LatitudeDecimal  float64
    LongitudeDecimal float64
    Time             time.Time

    // NOTE(dustin): !! Finish implementing.
    // Elevation      float32
    // Course         float32
    // Speed          float32
    // Hdop           float32
    // Src            string
    // SatelliteCount uint8
}

func (gts *GpxTrackSegmentBuilder) TrackPoint() *GpxTrackPointBuilder {
    return &GpxTrackPointBuilder{
        b: gts.b,
    }
}

func (gtpb *GpxTrackPointBuilder) Write() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if gtpb.Time.IsZero() {
        log.Panicf("timestamp not set")
    }

    // TODO(dustin): !! Handle epsilon.
    if gtpb.LatitudeDecimal == 0.0 {
        log.Panicf("latitude not set")
    }

    if gtpb.LongitudeDecimal == 0.0 {
        log.Panicf("longitude not set")
    }

    attrs := make([]xml.Attr, 2)
    attrs[0] = xml.Attr{Name: xml.Name{"", "lat"}, Value: strconv.FormatFloat(gtpb.LatitudeDecimal, 'f', -1, 64)}
    attrs[1] = xml.Attr{Name: xml.Name{"", "lon"}, Value: strconv.FormatFloat(gtpb.LongitudeDecimal, 'f', -1, 64)}

    trkptStart := xml.StartElement{
        Name: xml.Name{
            Space: "",
            Local: "trkpt",
        },
        Attr: attrs,
    }

    err = gtpb.b.encoder.EncodeToken(trkptStart)
    log.PanicIf(err)

    timeStart := xml.StartElement{
        Name: xml.Name{
            Space: "",
            Local: "time",
        },
    }

    err = gtpb.b.encoder.EncodeElement(gtpb.Time.UTC().Format("2006-01-02T15:04:05-0700"), timeStart)
    log.PanicIf(err)

    trkptEnd := xml.EndElement{
        Name: xml.Name{
            Space: "",
            Local: "trkpt",
        },
    }

    err = gtpb.b.encoder.EncodeToken(trkptEnd)
    log.PanicIf(err)

    return nil
}
