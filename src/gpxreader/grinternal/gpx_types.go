package grinternal

import (
    "fmt"
    "time"
)

/*

- Add structures for the following nodes:

    rte
    link
    copyright
    email
    author
    metadata
    bounds

- What does the MovingData struct from the original project represent?
- Additional reference: http://www.topografix.com/gpx_manual.asp#hdop

*/

type Gpx struct {
    Xmlns string
    Xsi string
    Version float32
    Creator string
    SchemaLocation string
    Time time.Time
}

type Track struct {
}

type TrackSegment struct {

}

type TrackPoint struct {
    LatitudeDecimal float64
    LongitudeDecimal float64
    Elevation float32
    Course float32
    Speed float32
    Hdop float32
    Src string
    SatelliteCount uint8
    Time time.Time
}

func (tp *TrackPoint) String() string {
    return fmt.Sprintf("TrackPoint<LAT=(%.8f) LON=(%.8f) ELV=(%f) CRS=(%f) SPD=(%f) HDOP=(%f) SRC=[%s] SAT=(%d) TIME=[%s]>", tp.LatitudeDecimal, tp.LongitudeDecimal, tp.Elevation, tp.Course, tp.Speed, tp.Hdop, tp.Src, tp.SatelliteCount, tp.Time)
}
