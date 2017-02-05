package gpxreader

import (
    "testing"
    "time"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-time-index"
)

func checkIndex(t *testing.T, gi *GpxIndex, label string, ti timeindex.TimeInterval, indexPosition int) {
    
    // `indexPosition` describes where the start-timestamp should be in the 
    // index.

    if _, found := gi.members[label]; found == false {
        t.Fatalf("Member not present.")
    } else if gi.fileTimes[indexPosition] != ti {
        t.Fatalf("Indexed start time is not correct.")
    }

    files, found := gi.files[ti]
    if found == false {
        t.Fatalf("File info not found.")
    }

    gfi := files[0]
    if gfi.Label != label {
        t.Fatalf("GFI label is not correct.")
    } else if gfi.lastPointTime != ti[1] {
        t.Fatalf("GFI stop time is not correct.")
    }
}

func getTestGpxIndexAccessor() *GpxBufferedDataAccessor {
    label1 := "testfile1"
    label2 := "testfile2"
    
    gbda := NewGpxBufferedDataAccessor()

    if err := gbda.Add(label1, TestGpxData); err != nil {
        log.Panic(err)
    }

    if err := gbda.Add(label2, TestGpxData2); err != nil {
        log.Panic(err)
    }

    return gbda
}

func TestGpxIndexAdd(t *testing.T) {
    label1 := "testfile1"
    label2 := "testfile2"
    
    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 10 * time.Minute, 0)

    // Note that we add TestGpxData2 before TestGpxData because the start-time 
    // of TestGpxData should be inserted before TestGpxData2. If something 
    // breaks down in the storage/sorting, there's a better chance that it'll 
    // stand-out this way.

    // Check that the first file is properly indexed.
    
    t2, err := gi.Add(label2)
    log.PanicIf(err)

    // Check that the second file is properly indexed.

    t1, err := gi.Add(label1)
    log.PanicIf(err)

    checkIndex(t, gi, label1, t1, 0)
    checkIndex(t, gi, label2, t2, 1)
}

func TestGpxIndexEnsureLoaded(t *testing.T) {
    label1 := "testfile1"

    // Load data.

    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 10 * time.Minute, 1)

    _, err := gi.Add(label1)
    log.PanicIf(err)

    if len(gi.mru) > 0 {
        t.Fatalf("MRU is non-empty before load.")
    }

    // Check GFI 1.

    gfi1 := gi.members[label1]

    if err := gi.ensureLoaded(gfi1); err != nil {
        log.Panic(err)
    }

    if len(gi.mru) != 1 {
        t.Fatalf("MRU does not have exactly one item after first add.")
    } else if gfi1.isLoaded != true {
        t.Fatalf("GFI 1 isn't marked as loaded.")
    } else if len(gfi1.index) != 204 {
        t.Fatalf("GFI 1 index has wrong size: (%d)", len(gfi1.index))
    } else if len(gfi1.points) != 204 {
        t.Fatalf("GFI 1 points has wrong size: (%d)", len(gfi1.points))
    }

    first := gfi1.index[0]
    last := gfi1.index[len(gfi1.index) - 1]

    if first.Format(time.RFC3339) != "2016-12-02T08:05:44Z" {
        t.Fatalf("First time in index is not correct.")
    } else if last.Format(time.RFC3339) != "2016-12-03T07:57:07Z" {
        t.Fatalf("Last time in index is not correct.")
    }

    left, err := time.Parse(time.RFC3339, "2016-12-02T16:16:23Z")
    log.PanicIf(err)

    right, err := time.Parse(time.RFC3339, "2016-12-02T16:27:08Z")
    log.PanicIf(err)

    q, err := time.Parse(time.RFC3339, "2016-12-02T16:23:23Z")
    log.PanicIf(err)

    results := getNearestTimes(gfi1.index, q, time.Minute * 5)
    if len(results) != 1 || results[0] != right {
        t.Fatalf("GFI 1 one-ended nearest-time search failed.")
    }

    results = getNearestTimes(gfi1.index, q, time.Minute * 8)
    if len(results) != 2 || results[0] != left || results[1] != right {
        t.Fatalf("GFI 1 two-ended nearest-time search failed.")
    }
}

func getNearestTimes(ts timeindex.TimeSlice, q time.Time, tolerance time.Duration) (results []time.Time) {
    results = make([]time.Time, 0)

    cb := func(t time.Time) error {
        results = append(results, t)
        return nil
    }

    if err := ts.SearchNearest(q, tolerance, cb); err != nil {
        log.Panic(err)
    }

    return results
}

func TestGpxIndexSearchExact(t *testing.T) {
    label1 := "testfile1"
    label2 := "testfile2"

    // Load data.

    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 1 * time.Minute, 1)

    _, err := gi.Add(label1)
    log.PanicIf(err)

    _, err = gi.Add(label2)
    log.PanicIf(err)

    // Test the autoload of data-set 1. Do an exact-match.

    q, err := time.Parse(time.RFC3339, "2016-12-03T07:23:50Z")
    log.PanicIf(err)

    matches, err := gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 1 {
        t.Fatalf("MRU not right size after search into first data-set (1): (%d)", len(gi.mru))
    } else if gi.mru[0] != label1 {
        t.Fatalf("MRU not populated with source.")
    }

    if len(matches) != 1 || matches[0].Time != q {
        t.Fatalf("GI exact search results not correct.")
    } 

    m := matches[0]

    if m.FileInfo.Label != label1 {
        t.Fatalf("Match label not correct.")
    } else if m.Point.Latitude != 48.45672936673762 {
        t.Fatalf("Match latitude not correct.")
    } else if m.Point.Longitude != -122.34140644128601 {
        t.Fatalf("Match longitude not correct.")
    }
}

func TestGpxIndexSearchApproximate(t *testing.T) {
    label1 := "testfile1"
    label2 := "testfile2"

    // Load data.

    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 4 * time.Minute, 1)

    _, err := gi.Add(label1)
    log.PanicIf(err)

    _, err = gi.Add(label2)
    log.PanicIf(err)

    q, err := time.Parse(time.RFC3339, "2016-12-03T07:26:00Z")
    log.PanicIf(err)

    matches, err := gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 1 {
        t.Fatalf("MRU not right size after search into first data-set (2).")
    } else if gi.mru[0] != label1 {
        t.Fatalf("MRU not populated with source.")
    }

    left, err := time.Parse(time.RFC3339, "2016-12-03T07:23:50Z")
    log.PanicIf(err)

    right, err := time.Parse(time.RFC3339, "2016-12-03T07:29:20Z")
    log.PanicIf(err)

    len1 := len(matches)
    if len1 != 2 {
        t.Fatalf("GI approximate search results size not correct: (%d)", len1)
    } else if matches[0].Time != left || matches[1].Time != right {
        t.Fatalf("GI approximate search results not correct.")
    }
}

func TestGpxIndexSearchSecondDatasetWithMru(t *testing.T) {
    label1 := "testfile1"
    label2 := "testfile2"

    // Load data.

    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 1 * time.Minute, 2)

    _, err := gi.Add(label1)
    log.PanicIf(err)

    _, err = gi.Add(label2)
    log.PanicIf(err)

    // Invoke the first dataset.

    q, err := time.Parse(time.RFC3339, "2016-12-03T07:26:00Z")
    log.PanicIf(err)

    matches, err := gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 1 {
        t.Fatalf("MRU not right size after search into first data-set.")
    } else if gi.mru[0] != label1 {
        t.Fatalf("MRU not populated with source.")
    }

    // Invoke the second data-set.

    q, err = time.Parse(time.RFC3339, "2016-12-22T07:13:21Z")
    log.PanicIf(err)

    matches, err = gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 2 {
        t.Fatalf("MRU not right size after search into second data-set.")
    } else if gi.mru[0] != label2 || gi.mru[1] != label1 {
        t.Fatalf("MRU not populated correctly.")
    }

    if len(matches) != 1 || matches[0].Time != q {
        t.Fatalf("GI exact search result from second data-set not correct.")
    }
}

func TestGpxIndexSearchSecondDatasetWithoutMru(t *testing.T) {
    label1 := "testfile1"
    label2 := "testfile2"

    // Load data.

    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 1 * time.Minute, 0)

    _, err := gi.Add(label1)
    log.PanicIf(err)

    _, err = gi.Add(label2)
    log.PanicIf(err)

    // Invoke the first dataset.

    q, err := time.Parse(time.RFC3339, "2016-12-03T07:26:00Z")
    log.PanicIf(err)

    matches, err := gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 0 {
        t.Fatalf("MRU not right size after search into first data-set.")
    }

    // Invoke the second data-set.

    q, err = time.Parse(time.RFC3339, "2016-12-22T07:13:21Z")
    log.PanicIf(err)

    matches, err = gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 0 {
        t.Fatalf("MRU not right size after search into second data-set.")
    }

    if len(matches) != 1 || matches[0].Time != q {
        t.Fatalf("GI exact search result from second data-set not correct.")
    }
}

func TestGpxIndexUnload(t *testing.T) {
    label1 := "testfile1"
    label2 := "testfile2"

    // Load data.

    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 1 * time.Minute, 1)

    _, err := gi.Add(label1)
    log.PanicIf(err)

    _, err = gi.Add(label2)
    log.PanicIf(err)

    // Invoke the first dataset.

    q, err := time.Parse(time.RFC3339, "2016-12-03T07:26:00Z")
    log.PanicIf(err)

    _, err = gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 1 {
        t.Fatalf("MRU not right size after search into first data-set.")
    } else if gi.mru[0] != label1 {
        t.Fatalf("MRU not populated with source.")
    }

    // Invoke the second data-set.

    // Make sure the first file is currently loaded.
    if gi.members[label1].isLoaded != true {
        t.Fatalf("First file is not already loaded.")
    }

    q, err = time.Parse(time.RFC3339, "2016-12-22T07:13:21Z")
    log.PanicIf(err)

    _, err = gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 1 {
        t.Fatalf("MRU not right size after search into second data-set.")
    } else if gi.mru[0] != label2 {
        t.Fatalf("MRU not populated correctly.")
    } else if gi.members[label1].isLoaded != false {
        t.Fatalf("First file was not supplanted by second file.")
    }
}

func TestGpxIndexMruPromote(t *testing.T) {
    label1 := "testfile1"
    label2 := "testfile2"

    // Load data.

    gbda := getTestGpxIndexAccessor()
    gi := NewGpxIndex(gbda, 1 * time.Minute, 3)

    _, err := gi.Add(label1)
    log.PanicIf(err)

    _, err = gi.Add(label2)
    log.PanicIf(err)

    // Invoke the first dataset.

    q, err := time.Parse(time.RFC3339, "2016-12-03T07:26:00Z")
    log.PanicIf(err)

    _, err = gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 1 {
        t.Fatalf("MRU not right size after search into first data-set.")
    } else if gi.mru[0] != label1 {
        t.Fatalf("MRU not populated with source.")
    }

    // Invoke the second data-set.

    q, err = time.Parse(time.RFC3339, "2016-12-22T07:13:21Z")
    log.PanicIf(err)

    _, err = gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 2 {
        t.Fatalf("MRU not right size after search into second data-set.")
    } else if gi.mru[0] != label2 || gi.mru[1] != label1 {
        t.Fatalf("MRU not populated correctly.")
    }

    // Invoke the first dataset.

    q, err = time.Parse(time.RFC3339, "2016-12-03T07:26:00Z")
    log.PanicIf(err)

    _, err = gi.Search(q)
    log.PanicIf(err)

    if len(gi.mru) != 2 {
        t.Fatalf("MRU not right size after search into first data-set.")
    } else if gi.mru[0] != label1 || gi.mru[1] != label2 {
        t.Fatalf("MRU not populated with source (promotion was expected).")
    }
}
