// XML parser/visitor logic

package xmlvisitor

import (
    "strings"
    "errors"
    "io"

    "encoding/xml"
)

const (
    XmlPart_Initial = iota
    XmlPartStartTag = iota
    XmlPartEndTag   = iota
    XmlPartCharData = iota
)

type SimpleXmlVisitor interface {
    // The content identifier next to the left angle brack.
    HandleStart(tagName *string, attrp *map[string]string, xp *XmlParser) error
    
    // The content identifier next to the left angle brack.
    HandleEnd(tagName *string, xp *XmlParser) error
    
    // Return the value that was found between two tags not having any child 
    // nodes.
    HandleValue(tagName *string, data *string, xp *XmlParser) error
}

type ExtendedXmlVisitor interface {
    // The content identifier next to the left angle brack.
    HandleStart(tagName *string, attrp *map[string]string, xp *XmlParser) error
    
    // The content identifier next to the left angle brack.
    HandleEnd(tagName *string, xp *XmlParser) error

    // Return the value that was found between two tags not having any child 
    // nodes.
    HandleValue(tagName *string, data *string, xp *XmlParser) error
    
    // Content that comes before one open/close tag and an adjacent one: either 
    // the useless whitespace between two open adjacent tags or two close 
    // adjacent tags or a tangible/empty value between an open and close tag.
    HandleCharData(data *string, xp *XmlParser) error

    // Example:
    //
    // <!-- Comment -->
    HandleComment(comment *string, xp *XmlParser) error
    
    // Example:
    //
    // <?xml version="1.0" encoding="UTF-8"?>
    HandleProcessingInstruction(target *string, instruction *string, xp *XmlParser) error
    
    // Example:
    //
    // <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
    //   "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
    //
    // <![CDATA[Some text here.]]>
    HandleDirective(directive *string, xp *XmlParser) error
}

type XmlVisitor interface{}

type XmlParser struct {
    decoder *xml.Decoder
    ns *Stack
    v XmlVisitor
    
    // The state preceding the last state.
    lastState2 int
    
    // The last state.
    lastState1 int

    lastCharData string
    
    doReportMarginCharData bool
    doAutoTrimCharData bool
}

func (xp *XmlParser) PushLastState(lastState int) {
    xp.lastState2 = xp.lastState1
    xp.lastState1 = lastState
}

func (xp *XmlParser) GetLastState() int {
    return xp.lastState1
}

func (xp *XmlParser) GetLastLastState() int {
    return xp.lastState2
}

func (xp *XmlParser) NodeStack() *Stack {
    return xp.ns
}

func (xp *XmlParser) SetDoReportMarginCharData(value bool) {
    xp.doReportMarginCharData = value
}

func (xp *XmlParser) SetDoAutoTrimCharData(value bool) {
    xp.doAutoTrimCharData = value
}

// Create parser.
func NewXmlParser(r io.Reader, visitor XmlVisitor) *XmlParser {
    decoder := xml.NewDecoder(r)
    ns := newStack()

    return &XmlParser {
            decoder: decoder,
            ns: ns,
            v: visitor,
            lastState1: XmlPart_Initial,
            lastState2: XmlPart_Initial,
            doReportMarginCharData: false,
            doAutoTrimCharData: true,
    }
}

func (xp *XmlParser) LastState() int {
    return xp.lastState1
}

func (xp *XmlParser) LastStateName() string {
    if xp.lastState1 == XmlPart_Initial {
        return ""
    } else if xp.lastState1 == XmlPartStartTag {
        return "StartTag"
    } else if xp.lastState1 == XmlPartEndTag {
        return "EndTag"
    } else if xp.lastState1 == XmlPartCharData {
        return "CharData"
    } else {
        panic(errors.New("Invalid XML state."))
    }
}

// Run the parse with a minimal memory footprint.
func (xp *XmlParser) Parse() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    for {
        token, err := xp.decoder.Token()
        if err != nil {
            break
        }
  
        switch e := token.(type) {
        case xml.StartElement:
            name := e.Name.Local

            xp.ns.push(name)

            var attributes map[string]string = make(map[string]string)
            for _, a := range e.Attr {
                attributes[a.Name.Local] = a.Value
            }

            sxv := xp.v.(SimpleXmlVisitor)
            err := sxv.HandleStart(&name, &attributes, xp)
            if err != nil {
                panic(err)
            }

            xp.PushLastState(XmlPartStartTag)

        case xml.EndElement:
            xp.ns.pop()

            name := e.Name.Local

            sxv := xp.v.(SimpleXmlVisitor)
            err := sxv.HandleEnd(&name, xp)
            if err != nil {
                panic(err)
            }

            if xp.lastState1 == XmlPartCharData && xp.lastState2 == XmlPartStartTag {
                sxv.HandleValue(&name, &xp.lastCharData, xp)
            }

            xp.PushLastState(XmlPartEndTag)

        case xml.CharData:
            // The underlying/aliased type is byte[].
            s := string(e)
            
            if xp.doAutoTrimCharData == true {
                s = strings.TrimSpace(s)
            }

            // We'll be reporting values found between open and close tags in 
            // the tag-end case. Therefore, this will be exclusively 
            // responsible for reporting the data between adjacent open tags or 
            // adjacent close tags.
            if xp.doReportMarginCharData == true {
                sxv, ok := xp.v.(ExtendedXmlVisitor)
                if ok == true {
                    err := sxv.HandleCharData(&s, xp)
                    if err != nil {
                        panic(err)
                    }
                }
            }

            xp.lastCharData = s
            xp.PushLastState(XmlPartCharData)

        case xml.Comment:
            // The underlying/aliased type is byte[].
            s := string(e)

            exv, ok := xp.v.(ExtendedXmlVisitor)
            if ok == true {
                err := exv.HandleComment(&s, xp)
                if err != nil {
                    panic(err)
                }
            }

        case xml.ProcInst:
            instruction := string(e.Inst)

            exv, ok := xp.v.(ExtendedXmlVisitor)
            if ok == true {
                err := exv.HandleProcessingInstruction(&e.Target, &instruction, xp)
                if err != nil {
                    panic(err)
                }
            }

        case xml.Directive:
            // The underlying/aliased type is byte[].
            s := string(e)

            exv, ok := xp.v.(ExtendedXmlVisitor)
            if ok == true {
                err := exv.HandleDirective(&s, xp)
                if err != nil {
                    panic(err)
                }
            }
        }
    }

    return nil
}
