package grinternal

import (
    "strconv"
    "time"
)

func parseFloat32(raw string) float32 {
    v, err := strconv.ParseFloat(raw, 32)
    if err != nil {
        panic(err)
    }

    return float32(v)
}

func parseFloat64(raw string) float64 {
    v, err := strconv.ParseFloat(raw, 64)
    if err != nil {
        panic(err)
    }

    return v
}

func parseUint8(raw string) uint8 {
    v, err := strconv.ParseUint(raw, 10, 8)
    if err != nil {
        panic(err)
    }

    return uint8(v)
}

func parseIso8601Time(raw string) time.Time {
    t, err := time.Parse(time.RFC3339Nano, raw)
    if err != nil {
        panic(err)
    }

    return t
}
