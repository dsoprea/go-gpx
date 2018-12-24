## Description

The native Go XML package wants to load the whole XML document into memory. It does allow you to incrementally traverse the nodes, but, the moment that you tell it to decode one, all child nodes will be discarded. This means that wanting to look at the base GPX info in the root node comes at a cost of discarding all of the recorded points several levels below it.

This project uses the basic XML tokenization that the XML package provides while parsing the data and assigning the attributes and character-data to data-structures itself. It is very efficient in that no unnecessary seeking is done and no substantial amount of data is kept in memory. You simply provide a class that fulfills interface describing a bunch of callbacks and it is triggered at the various nodes with the information about that node/entity (persuant to the visitor pattern).


## Usage

The visitor type can satisfy the following interfaces (types are found in the `gpxreader` package):

```golang
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
```

Example usage, given a string with the GPX text. `NewGpsPointCollector` is a type that satisfies all of the interfaces. This is based on similar code from the tests:

```golang
b := bytes.NewBufferString(testGpxData)
gpc := NewGpsPointCollector()
gp := gpxreader.NewGpxParser(b, gpc)

if err := gp.Parse(); err != nil {
    panic(err)
}
```

There are also two convenience functions provided that allow you to avoid having to deal with interfaces if you only care about reading points:

```golang
// func EnumerateTrackPoints(f io.Reader, tpc TrackPointCallback) (err error)
// func ExtractTrackPoints(f io.Reader) (points []TrackPoint, err error)
```

`TrackPointCallback` is aliased to `func(tp *TrackPoint) error`.


## Indexing

We also provide the `GpxIndex` type to search for timestamps over a set of GPX files. Files are loaded on-demand. You can also specify a limit on the number of files loaded concurrently at any given time. A type that fulfills `GpxDataAccessor` must be provided in order to retrieve the GPX data.

Example:

```go
package main

import (
    "time"
    "fmt"

    "github.com/dsoprea/go-gpx/reader"
)

func findAndPrint(gi *gpxreader.GpxIndex, timePhrase string) {
    q, err := time.Parse(time.RFC3339, timePhrase)
    if err != nil {
        panic(err)
    }

    // Returns sorted first by time and then file-path.
    matches, err := gi.Search(q)
    if err != nil {
        panic(err)
    }

    for _, match := range matches {
        fmt.Printf("MATCH: [%s] (%f, %f) IN [%s]\n", match.Time, match.Point.Latitude, match.Point.Longitude, match.FileInfo.Label)
    }
}

func main() {
    tolerance := 5 * time.Minute
    maxFilesLoaded := 0

    gfda := new(gpxreader.GpxFileDataAccessor)
    gi := gpxreader.NewGpxIndex(gfda, tolerance, maxFilesLoaded)

    // Add() will return a `timeindex.TimeInterval` (`[2]time.Time`) that describes the range of time represented by the file.

    if _, err := gi.Add("trip_day1.gpx"); err != nil {
        panic(err)
    }

    if _, err := gi.Add("trip_day2.gpx"); err != nil {
        panic(err)
    }

    findAndPrint(gi, "2016-12-22T14:32:59Z")
    findAndPrint(gi, "2016-12-02T13:02:01Z")

    // Output:
    //
    // MATCH: [2016-12-22 14:30:59 +0000 UTC] (8.967136, -79.533077) IN [trip_day2.gpx]
    // MATCH: [2016-12-22 14:37:21 +0000 UTC] (8.967136, -79.533077) IN [trip_day2.gpx]
    // MATCH: [2016-12-02 13:00:01 +0000 UTC] (47.613163, -122.340196) IN [trip_day1.gpx]
    // MATCH: [2016-12-02 13:06:02 +0000 UTC] (47.613163, -122.340196) IN [trip_day1.gpx]
}
```


## To Do

- Only the primary location information and other data that is highly common and related is read and supported by the implemented types. We still need to add any attributes or parse any nodes that are currently missing per the spec (see the TODO file). Feel free to request this and/or submit a PR.
