## Description

The native Go XML package wants to load the whole XML document into memory. It does allow you to incremental traverse the nodes, but, the moment that you tell it to decode one, all child nodes will be discarded. This means that wanting to look at the base GPX info in the root node comes at a cost of discarding all of the recorded points several levels below it.

This project uses the basic XML tokenization that the XML package provides while parsing the data and assigning the attributes and character-data to data-structures itself. It is very efficient in that no unnecessary seeking is done and no substantial amount of data is kept in memory. You simply provide a class that fulfills interface describing a bunch of callbacks and it is triggered at the various nodes with the information about that node/entity (persuant to the visitor pattern).


## Example

The `gpxreadertest` tool is provided for a reference implementation of the `gpxreader` package (which is also the name of this library's package). 

To install:

```
$ go get github.com/dsoprea/go-gpxreader/commands/gpxreadertest
```

Most of the source:

```go
//...

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

//...

func main() {
    var gpxFilepath string

    o := readOptions()

    gpxFilepath = o.GpxFilepath

    f, err := os.Open(gpxFilepath)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    gv := newGpxVisitor()
    gp := gpxreader.NewGpxParser(f, gv)

    err = gp.Parse()
    if err != nil {
        print("Error: %s\n", err.Error())
        os.Exit(1)
    }
}
```

Output:

```
$ gpxreadertest -f 20140909.gpx 
GPX: GPX<C=[GPSLogger - http://gpslogger.mendhak.com/]>
Track: Track<>
Track segment: TrackSegment<>
Point: TrackPoint<LAT=(26.47886514) LON=(-80.08643986) ELV=(-12.000000) CRS=(197.899994) SPD=(35.250000) HDOP=(0.900000) SRC=[gps] SAT=(21) TIME=[2014-09-09 19:07:27 +0000 UTC]>
Point: TrackPoint<LAT=(26.40728154) LON=(-80.11801469) ELV=(9.000000) CRS=(0.000000) SPD=(0.000000) HDOP=(1.200000) SRC=[gps] SAT=(16) TIME=[2014-09-09 22:07:52 +0000 UTC]>
Point: TrackPoint<LAT=(26.54074478) LON=(-80.07230151) ELV=(-31.000000) CRS=(12.800000) SPD=(31.503967) HDOP=(1.000000) SRC=[gps] SAT=(17) TIME=[2014-09-09 22:53:27 +0000 UTC]>
```


## To Do

- Only the primary location information and other data that is highly common and related is read and supported by the provided data structures. We need to add any attributes or parse any nodes that are currently missing per the spec (see the TODO file). Feel free to request this and/or submit a PR.
