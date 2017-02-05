package main

import (
    "time"
    "fmt"

    "github.com/dsoprea/go-gpxreader"
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
