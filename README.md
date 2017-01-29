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


## To Do

- Only the primary location information and other data that is highly common and related is read and supported by the provided data structures. We need to add any attributes or parse any nodes that are currently missing per the spec (see the TODO file). Feel free to request this and/or submit a PR.
