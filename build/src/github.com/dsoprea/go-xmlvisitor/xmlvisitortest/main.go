package main

import (
    "os"
    "fmt"
    "strings"
    "io"

    "github.com/dsoprea/go-xmlvisitor"

    flags "github.com/jessevdk/go-flags"
)

type xmlVisitor struct {
}

func (xv *xmlVisitor) HandleStart(tagName *string, attrp *map[string]string, xp *xmlvisitor.XmlParser) error {
    fmt.Printf("Start: [%s]\n", *tagName)

    return nil
}

func (xv *xmlVisitor) HandleEnd(tagName *string, xp *xmlvisitor.XmlParser) error {
    fmt.Printf("Stop: [%s]\n", *tagName)

    return nil
}

func (xv *xmlVisitor) HandleValue(tagName *string, value *string, xp *xmlvisitor.XmlParser) error {
    fmt.Printf("Value: [%s] [%s]\n", *tagName, *value)

    return nil
}

/*
func (xv *xmlVisitor) HandleComment(comment *string, xp *xmlvisitor.XmlParser) error {
    fmt.Printf("Comment: [%s]\n", *comment)

    return nil
}

func (xv *xmlVisitor) HandleProcessingInstruction(target *string, instruction *string, xp *xmlvisitor.XmlParser) error {
    fmt.Printf("Processing Instruction: [%s] [%s]\n", *target, *instruction)

    return nil
}

func (xv *xmlVisitor) HandleDirective(directive *string, xp *xmlvisitor.XmlParser) error {
    fmt.Printf("Directive: [%s]\n", *directive)

    return nil
}
*/

func newXmlVisitor() (*xmlVisitor) {
    return &xmlVisitor {}
}

type options struct {
    XmlFilepath string  `short:"f" long:"xml-filepath" description:"XML file-path" required:"true"`
}

func readOptions () *options {
    o := options {}

    _, err := flags.Parse(&o)
    if err != nil {
        os.Exit(1)
    }

    return &o
}

func getTextReader() io.Reader {
    s := `<node1>
    <node2>
        <node3>node3 value</node3>
        <node4>node4 value</node4>
    </node2>
</node1>`

    r := strings.NewReader(s)

    return r
}

func getFileReader() io.Reader {
    var xmlFilepath string

    o := readOptions()
    xmlFilepath = o.XmlFilepath

    f, err := os.Open(xmlFilepath)
    if err != nil {
        panic(err)
    }

    return f
}

func closeFileReader(f os.File) {
    f.Close()
}

func main() {
//    r := getTextReader()

    r := getFileReader()
    f := r.(*os.File)
    defer closeFileReader(*f)

    v := newXmlVisitor()
    p := xmlvisitor.NewXmlParser(r, v)

    err := p.Parse()
    if err != nil {
        print("Error: %s\n", err.Error())
        os.Exit(1)
    }
}
