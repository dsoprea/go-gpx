package xmlvisitor

import (
)

type stackNode struct {
    value interface{}
    parent *stackNode
    next *stackNode
}

type Stack struct {
    root *stackNode
    last *stackNode
}

func newStack() *Stack {
    return &Stack {}
}

func (s *Stack) push(value interface{}) {
    newNode := &stackNode {
            parent: s.last,
            value: value,
    }

    if s.root == nil {
        s.root = newNode
    } else {
        s.last.next = newNode
    }

    s.last = newNode
}

func (s *Stack) pop() interface{} {
    if s.root == nil {
        return nil
    }

    lastValue := s.last.value
    s.last = s.last.parent

    // If the last node had a nil parent, it was the root.
    if s.last == nil {
        s.root = nil
    }

    return lastValue
}

// Return the value of the (n-i) node (n being the last node).
func (s *Stack) PeekFromEnd(i int) interface{} {
    n := i
    ptr := s.last

    for ptr != nil && n > 0 {
        ptr = ptr.parent
        n--
    }

    if ptr == nil {
        return nil
    }

    return ptr.value
}
