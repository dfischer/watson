// Package vm contains a stack-based virtual machine that interprets Watson.
package vm

import (
	"fmt"
)

const (
	DefaultStackSize = 1024 // the default size of the stack
)

// VM is a virtual machine that consists of a stack of values and a pointer to the top of the stack.
//
// In this document, we call the N-th element of the stack s[N] (counting from the top).
type VM struct {
	stack []*Value
	sp    int
}

// Returns a new VM with its stack allocated.
func NewVM() *VM {
	return &VM{
		stack: make([]*Value, DefaultStackSize),
		sp:    -1,
	}
}

// Op is an instruction executed by VM. Each op just manipulates the stack.
type Op int

const (
	// Integer Operations
	Inew Op = iota // push(0);
	Iinc           // v: int = pop(); push(v + 1);
	Ishl           // v: int = pop(); push(v << 1);
	Iadd           // b: int = pop(); a: int = pop(); push(a + b);

	// String Operations
	Snew // push("");
	Sadd // n: int = pop(); s: str = pop(); c = n && 0xff; push(s + c);

	// Object Operations
	Onew // push({});
	Oadd // v: any = pop(); k: str = pop(); o: obj = pop(); o[k] = v; push(o);

	// Bool Operations
	Bnew // push(false);
	Bneg // b: bool = pop(); push(!b);

	// Nil Operations
	Nnew // push(nil);
)

func (op Op) GoString() string {
	switch op {
	case Inew:
		return "Inew"
	case Iinc:
		return "Iinc"
	case Ishl:
		return "Ishl"
	case Iadd:
		return "Iadd"
	case Snew:
		return "Snew"
	case Sadd:
		return "Sadd"
	case Onew:
		return "Onew"
	case Oadd:
		return "Oadd"
	case Bnew:
		return "Bnew"
	case Bneg:
		return "Bneg"
	case Nnew:
		return "Nnew"
	default:
		panic(fmt.Errorf("invalid opcode: %d", op))
	}
}

var _ fmt.GoStringer = Op(0)

// Kind is a type of Value.
type Kind int

const (
	KInt    Kind = iota // 64-bit signed integer
	KString             // string (represented as a byte array)
	KObject             // object (set of key-value pairs)
	KBool               // bool
	KNil                // nil
)

// Value is an element of the stack.
type Value struct {
	Kind   Kind
	Int    int64
	String []byte
	Object Object
	Bool   bool
}

// Object is a set of key-value pairs.
type Object map[string]*Value

// NewIntValue creates a new Value that contains an integer.
func NewIntValue(val int64) *Value {
	return &Value{Kind: KInt, Int: val}
}

// NewStringValue creates a new Value that contains a string.
func NewStringValue(val []byte) *Value {
	return &Value{Kind: KString, String: val}
}

// NewObjectValue creates a new Value that contains an object.
func NewObjectValue(val Object) *Value {
	return &Value{Kind: KObject, Object: val}
}

// NewBoolValue creates a new Value that contains a bool.
func NewBoolValue(val bool) *Value {
	return &Value{Kind: KBool, Bool: val}
}

// NewNilValue creates a new Value that contains nil.
func NewNilValue() *Value {
	return &Value{Kind: KNil}
}

func (v *Value) DeepCopy() *Value {
	clone := &Value{Kind: v.Kind}
	switch v.Kind {
	case KInt:
		clone.Int = v.Int
	case KString:
		clone.String = make([]byte, len(v.String))
		copy(clone.String, v.String)
	case KObject:
		clone.Object = map[string]*Value{}
		for k, v := range v.Object {
			clone.Object[k] = v.DeepCopy()
		}
	case KBool:
		clone.Bool = v.Bool
	case KNil:
		// nop
	default:
		panic(fmt.Errorf("unknown kind: %d", v.Kind))
	}
	return clone
}
