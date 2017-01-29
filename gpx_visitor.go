package gpxreader

import (
    "io"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-xmlvisitor/xmlvisitor"
)


type GpxParser struct {
    xp *xmlvisitor.XmlParser
}

// Create parser.
func NewGpxParser(r io.Reader, visitor interface{}) *GpxParser {
    gp := new(GpxParser)

    v := newXmlVisitor(gp, visitor)
    gp.xp = xmlvisitor.NewXmlParser(r, v)

    return gp
}

// Run the parse with a minimal memory footprint.
func (gp *GpxParser) Parse() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if err := gp.xp.Parse(); err != nil {
        log.Panic(err)
    }

    return nil
}
