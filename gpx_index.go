// Provides a tool to dynamically index and provide lookups against a set of physical GPX files.
package gpxreader

import (
    "time"
    "fmt"
    "strings"
    "io"
    "os"
    "sort"
    "bytes"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-time-index"
)

var (
    ErrFileAlreadyAdded = fmt.Errorf("file already added")
    ErrEmptyFile = fmt.Errorf("file empty")
    ErrNotFound = fmt.Errorf("not found")
)


// GpxPoint A single location represented in a GPX file.
type GpxPoint struct {
    Latitude, Longitude float64
}


// GpxFileInfo Represents a single loaded GPX file.
type GpxFileInfo struct {
    // label Name or filepath of the GPX file.
    Label string

    // lastPointTime latest timestamp represented by the file.
    lastPointTime time.Time

    // count Number of points in file.
    count int

    // isLoaded Whether the points are currently loaded.
    isLoaded bool

    // index All of the times for points represented by the file.
    index timeindex.TimeSlice

    // points All points represented by the file keyed by time.
    points map[time.Time]GpxPoint
}


type GpxDataAccessor interface {
    Accessor(label string) (rc io.ReadCloser, err error)
}


type GpxFileDataAccessor struct {

}

func (gfda *GpxFileDataAccessor) Accessor(filepath string) (f io.ReadCloser, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    f, err = os.Open(filepath)
    log.PanicIf(err)

    return f, nil
}


type GpxBufferedDataAccessorResource struct {
    b *bytes.Buffer
}

func NewGpxBufferedDataAccessorResource(data string) *GpxBufferedDataAccessorResource {
    b := bytes.NewBufferString(data)

    return &GpxBufferedDataAccessorResource{
        b: b,
    }
}

func (gbdar GpxBufferedDataAccessorResource) Read(p []byte) (n int, err error) {
    n, err = gbdar.b.Read(p)
    return n, err
}

func (gbdar GpxBufferedDataAccessorResource) Close() (err error) {
    return nil
}


type GpxBufferedDataAccessor struct {
    sources map[string]string
}

func NewGpxBufferedDataAccessor() *GpxBufferedDataAccessor {
    sources := make(map[string]string)

    return &GpxBufferedDataAccessor{
        sources: sources,
    }
}

func (gbda *GpxBufferedDataAccessor) Add(label string, data string) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if _, found := gbda.sources[label]; found == true {
        log.Panic(fmt.Errorf("label already added"))
    }

    gbda.sources[label] = data

    return nil
}

func (gbda *GpxBufferedDataAccessor) Accessor(label string) (rc io.ReadCloser, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    data, found := gbda.sources[label]
    if found == false {
        log.Panic(fmt.Errorf("label not found"))
    }

    b := NewGpxBufferedDataAccessorResource(data)
    return b, nil
}

// GpxIndex Knows how to find location for a given point in time.
type GpxIndex struct {
    // gda Knows how to get file-data.
    gda GpxDataAccessor

    // fileTimes A sorted slice of the earliest represented times found in each 
    // file.
    fileTimes timeindex.TimeIntervalSlice

    // files A lookup of the earliest represent times found and the list of 
    // files that had them (we allow for the same earliest point to occur in 
    // more than one file).
    files map[timeindex.TimeInterval][]*GpxFileInfo
    
    // members A lookup of all known files (so we don't allow to load more than 
    // once).
    members map[string]*GpxFileInfo

    // proximityTolerance How near the previous point in the index has to be 
    // from the time that was queried to be allowed as a match. 
    proximityTolerance time.Duration

    // maxOpenFiles Maximum allowed open files.
    maxOpenFiles int

    // mru List of labels of open files sorted by usage (only if 
    // maxOpenFiles is greater than zero).
    mru []string
}

func NewGpxIndex(gda GpxDataAccessor, proximityTolerance time.Duration, maxOpenFiles int) *GpxIndex {
    ft := make(timeindex.TimeIntervalSlice, 0)
    f := make(map[timeindex.TimeInterval][]*GpxFileInfo)
    m := make(map[string]*GpxFileInfo)
    mru := make(sort.StringSlice, 0)

    return &GpxIndex{
        gda: gda,
        fileTimes: ft,
        files: f,
        members: m,
        proximityTolerance: proximityTolerance,
        maxOpenFiles: maxOpenFiles,
        mru: mru,
    }
}

func (gi *GpxIndex) Add(label string) (ti timeindex.TimeInterval, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    a, err := gi.gda.Accessor(label)
    log.PanicIf(err)

    defer a.Close()

    // Determine if this file has already been loaded.

    label = strings.ToLower(label)

    if _, found := gi.members[label]; found == true {
        return timeindex.TimeInterval{}, ErrFileAlreadyAdded
    }

    // Read the file once to establish the time range.

    gs, err := Summary(a)
    log.PanicIf(err)

    if gs.Count == 0 {
        return timeindex.TimeInterval{}, ErrEmptyFile
    }

    // Add to the list of start times representing all known files.

    ti = timeindex.TimeInterval { gs.Start, gs.Stop }
    gi.fileTimes = gi.fileTimes.Add(ti)

    // Store the file info.

    gfi := &GpxFileInfo{
        Label: label,
        lastPointTime: gs.Stop,
        count: gs.Count,
        index: make(timeindex.TimeSlice, 0),
        points: make(map[time.Time]GpxPoint),
    }

    gi.members[label] = gfi

    if files, found := gi.files[ti]; found == true {
        gi.files[ti] = append(files, gfi)
    } else {
        gi.files[ti] = []*GpxFileInfo { gfi }
    }

    return ti, nil
}

func (gi *GpxIndex) AddWithFile(filepath string) (ti timeindex.TimeInterval, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    f, err := os.Open(filepath)
    log.PanicIf(err)

    defer f.Close()

    ti, err = gi.Add(filepath)
    log.PanicIf(err)

    return ti, err
}


type IndexHits struct {
    Time time.Time
    Point GpxPoint
    FileInfo *GpxFileInfo
}

type IndexHitSlice []IndexHits

func (ihs IndexHitSlice) Search(ih IndexHits) int {
    return SearchIndexHits(ihs, ih)
}

func (ihs IndexHitSlice) Add(ih IndexHits) (newIhs IndexHitSlice) {
    i := ihs.Search(ih)
    if i < len(ihs) && ihs[i].Time == ih.Time {
        return
    }

    right := append(IndexHitSlice { ih }, ihs[i:]...)
    newIhs = append(ihs[:i], right...)

    return newIhs
}

func SearchIndexHits(ihs IndexHitSlice, ih IndexHits) int {
    p := func(i int) bool {
        return ihs[i].Time.After(ih.Time) || ihs[i].Time == ih.Time && ihs[i].FileInfo.Label > ih.FileInfo.Label
    }
    
    return Search(len(ihs), p)
}

func Search(n int, f func(int) bool) int {
    // Define f(-1) == false and f(n) == true.
    // Invariant: f(i-1) == false, f(j) == true.
    i, j := 0, n
    for i < j {
        h := i + (j-i)/2 // avoid overflow when computing h
        // i ≤ h < j
        if !f(h) {
            i = h + 1 // preserves f(i-1) == false
        } else {
            j = h // preserves f(j) == true
        }
    }

    // i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
    return i
}

// ensureLoaded Make sure the given GPX data is loaded. Unload older data if we 
// have to make room.
func (gi *GpxIndex) ensureLoaded(gfi *GpxFileInfo) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    // Promote this entry in the MRU if we're using one.
    if gfi.isLoaded == true && gi.maxOpenFiles > 0 {
        i := -1
        for j, label := range gi.mru {
            if label == gfi.Label {
                i = j
                break
            }
        }

        if i == -1 {
            log.Panic(fmt.Errorf("Could not found loaded file in MRU: [%s]", gfi.Label))
        }

        right := append(gi.mru[:i], gi.mru[i + 1:]...)
        gi.mru = append([]string { gfi.Label }, right...)

        return nil
    }

    // We're maxed-out on the files we've loaded. Deallocate the point data on 
    // the least-used, currently-allocated data.
    len_ := len(gi.mru)
    if gi.maxOpenFiles > 0 && len_ >= gi.maxOpenFiles {
        oldestFilepath := gi.mru[len_ - 1]
        oldestGfi := gi.members[oldestFilepath]

        oldestGfi.index = nil
        oldestGfi.points = nil
        oldestGfi.isLoaded = false

        gi.mru = gi.mru[:len_ - 1]
    }

    // Load the points.

    gfi.index = make(timeindex.TimeSlice, 0)
    gfi.points = make(map[time.Time]GpxPoint)

    tpc := func(tp *TrackPoint) (err error) {
        gfi.index = gfi.index.Add(tp.Time)

        gfi.points[tp.Time] = GpxPoint{
            Latitude: tp.LatitudeDecimal,
            Longitude: tp.LongitudeDecimal,
        }

        return nil
    }

    a, err := gi.gda.Accessor(gfi.Label)
    log.PanicIf(err)

    defer a.Close()

    if err := EnumerateTrackPoints(a, tpc); err != nil {
        log.Panic(err)
    }

    // Update our tracking information.

    gfi.isLoaded = true

    if gi.maxOpenFiles > 0 {
        gi.mru = append([]string { gfi.Label }, gi.mru...)
    }

    return nil
}

func (gi *GpxIndex) searchFile(matches []IndexHits, fileInterval timeindex.TimeInterval, searchTime time.Time) (newMatches []IndexHits, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    files := gi.files[fileInterval]
    results := IndexHitSlice(matches)
    for _, gfi := range files {
        if err := gi.ensureLoaded(gfi); err != nil {
            log.Panic(err)
        }

        cb := func(foundTime time.Time) (err error) {
            ih := IndexHits{
                Time: foundTime,
                Point: gfi.points[foundTime],
                FileInfo: gfi,
            }

            results = results.Add(ih)
            return nil
        }

        if err := gfi.index.SearchNearest(searchTime, gi.proximityTolerance, cb); err != nil {
            log.Panic(err)
        }
    }

    newMatches = []IndexHits(results)
    return newMatches, nil
}

// Search Determine if we have a point with a near-enough timestamp.
func (gi *GpxIndex) Search(t time.Time) (matches []IndexHits, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if len(gi.fileTimes) == 0 {
        return nil, ErrNotFound
    }

    matches = make([]IndexHits, 0)

    cb := func(ti timeindex.TimeInterval) (err error) {
        defer func() {
            if state := recover(); state != nil {
                err = log.Wrap(state.(error))
            }
        }()

        matches, err = gi.searchFile(matches, ti, t)
        log.PanicIf(err)

        return nil
    }

    if err := gi.fileTimes.Search(t, cb); err != nil {
        log.Panic(err)
    }

    return matches, nil
}
