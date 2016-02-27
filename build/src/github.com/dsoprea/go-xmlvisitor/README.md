## Introduction

This library allows for efficient traversal of XML using the visitor pattern within Go. It also tracks state in order to provide some convenience functions.

## Usage

Implement an interface and pass a reader resource along with an instance of your type to the parser.


### Interfaces

#### SimpleXmlVisitor

- `HandleStart(tagName *string, attrp *map[string]string, xp *XmlParser) error`

  Receives start nodes.

- `HandleEnd(tagName *string, xp *XmlParser) error`

  Receives end nodes.

- `HandleValue(tagName *string, data *string, xp *XmlParser) error`

  Return values found between start and end nodes where other subnodes weren't present. This is a convenience wrapper that intelligently manages character-data.

#### ExtendedXmlVisitor

- `HandleStart(tagName *string, attrp *map[string]string, xp *XmlParser) error`

  Receives start nodes.

- `HandleEnd(tagName *string, xp *XmlParser) error`

  Receives end nodes.

- `HandleValue(tagName *string, data *string, xp *XmlParser) error`

  Receives values found between start and end nodes where other subnodes weren't present.

- `HandleCharData(data *string, xp *XmlParser) error`

  Receives content ("character data") not found within a tag.

- `HandleComment(comment *string, xp *XmlParser) error`

  Receives comment text.

- `HandleProcessingInstruction(target *string, instruction *string, xp *XmlParser) error`

  Receives processing instructions (e.g. <?xml version="1.0" encoding="UTF-8"?>).

- `HandleDirective(directive *string, xp *XmlParser) error`

  Receives directives (e.g. <![CDATA[Some text here.]]>).


### Configuration

- SetDoReportMarginCharData(value bool)

  Default: false

  Trigger on the character data that appears between adjacent open tags or 
  adjacent close tags.

- SetDoAutoTrimCharData(value bool)

  Default: true

  Remove empty space from the ends of character data. This also affects the 
  values that we derive from character data (received by HandleValue()).


## Other Conveniences

The visitor callbacks will have an instance of the XmlParser passed-in as an argument. This can be used to access additional functionality.

### Last Node State

Calling `GetLastState()` on the XmlParser object will return the last [useful] type of token that was encountered. It will be equal to one of the *XmlPart* constants:

- `xmlvisitor.XmlPartStartTag`
- `xmlvisitor.XmlPartEndTag`
- `xmlvisitor.XmlPartCharData`


### Stack

The visitor callbacks have access to the current stack of nodes using `NodeStack()`. This returns an instance of the *Stack* type. See [stack.go](src/gpxreader/gpxreader/stack.go) for further detail.


## Notes

- There is no specific handling of namespaces. This is left as an exercise to the implementor. This library merely provides a simplification of the low-level tokenizer.


## Example

```go
package main

import (
    "os"
    "fmt"
    "strings"
    "io"

    "github.com/dsoprea/go-xmlvisitor"
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

func newXmlVisitor() (*xmlVisitor) {
    return &xmlVisitor {}
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

func main() {
    r := getTextReader()

    v := newXmlVisitor()
    p := xmlvisitor.NewXmlParser(r, v)

    err := p.Parse()
    if err != nil {
        print("Error: %s\n", err.Error())
        os.Exit(1)
    }
}
```

Output:

```
Start: [node1]
Start: [node2]
Start: [node3]
Stop: [node3]
Value: [node3] [node3 value]
Start: [node4]
Stop: [node4]
Value: [node4] [node4 value]
Stop: [node2]
Stop: [node1]
```
