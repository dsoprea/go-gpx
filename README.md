## Description

The native Go XML package wants to load the whole XML document into memory. It does allow you to browse the nodes, but the moment that you tell it to decode it, it will discard all subnodes. This means that wanting to look at the base GPX info in the root node comes at a cost of discarding all of the recorded points several levels below it.

This project uses the basic XML tokenization that the XML package provides while parsing the data and assigning the attributes and character-data to data-structures itself. It is very efficient in that no data is processed before you're ready for it and no substantial amount of data is kept in memory. You simply provide a class that fulfills an callback interface and it's triggered at the various nodes with the information about that node/entity.


## Example

The `gpxparse` tool and a modest GPS log is provided for a reference implementation of the `gpxreader` package (which is also the name of this library's package). This is the mostly the source of that tool:

```go
package main

import (
    "os"
    "fmt"

    "gpxreader/gpxreader"
)

type gpxVisitor struct {}

func newgpxVisitor() (*gpxVisitor) {
    return &gpxVisitor {}
}

func (gv *gpxVisitor) GpxOpen(gpx *gpxreader.Gpx) error {
    fmt.Printf("GPX: %s\n", gpx)

    return nil
}

func (gv *gpxVisitor) GpxClose(gpx *gpxreader.Gpx) error {
    return nil
}

func (gv *gpxVisitor) TrackOpen(track *gpxreader.Track) error {
    fmt.Printf("Track: %s\n", track)

    return nil
}

func (gv *gpxVisitor) TrackClose(track *gpxreader.Track) error {
    return nil
}

func (gv *gpxVisitor) TrackSegmentOpen(trackSegment *gpxreader.TrackSegment) error {
    fmt.Printf("Track segment: %s\n", trackSegment)

    return nil
}

func (gv *gpxVisitor) TrackSegmentClose(trackSegment *gpxreader.TrackSegment) error {
    return nil
}

func (gv *gpxVisitor) TrackPointOpen(trackPoint *gpxreader.TrackPoint) error {
    return nil
}

func (gv *gpxVisitor) TrackPointClose(trackPoint *gpxreader.TrackPoint) error {
    fmt.Printf("Point: %s\n", trackPoint)

    return nil
}

func main() {
    var gpxFilepath string = "testdata/20130729.gpx"

    gv := newgpxVisitor()
    gp := gpxreader.NewGpxParser(&gpxFilepath, gv)

    defer gp.Close()

    err := gp.Parse()
    if err != nil {
        print("Error: %s\n", err.Error())
        os.Exit(1)
    }
}
```

Output:

```
$ bin/gpxparse -f test/20130729.gpx 
GPX: GPX<C=[GPSLogger - http://gpslogger.mendhak.com/]>
Track: Track<>
Track segment: TrackSegment<>
Point: TrackPoint<LAT=(26.07072655) LON=(-80.14360848) ELV=(-18.500000) CRS=(0.000000) SPD=(0.750000) HDOP=(5.800000) SRC=[gps] SAT=(7) TIME=[2013-07-30 02:38:29 +0000 UTC]>
Point: TrackPoint<LAT=(26.07099936) LON=(-80.14324075) ELV=(-43.900002) CRS=(0.000000) SPD=(0.000000) HDOP=(47.599998) SRC=[gps] SAT=(4) TIME=[2013-07-30 02:39:15 +0000 UTC]>
Point: TrackPoint<LAT=(26.07173904) LON=(-80.14322448) ELV=(-8.100000) CRS=(0.000000) SPD=(0.000000) HDOP=(23.600000) SRC=[gps] SAT=(5) TIME=[2013-07-30 02:40:17 +0000 UTC]>
Point: TrackPoint<LAT=(26.07182345) LON=(-80.14294142) ELV=(-24.600000) CRS=(0.000000) SPD=(0.000000) HDOP=(21.600000) SRC=[gps] SAT=(5) TIME=[2013-07-30 02:41:19 +0000 UTC]>
```


## To Do

- Only the primary location information and other data that is highly common and related is read and supported by the provided data structures. We need to add any attributes or parse any nodes that are currently missing per the spec (see the TODO file). Feel free to request this and/or submit a PR.
