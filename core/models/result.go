package models

import "bytes"

type Code int32

const (
	UNSPECIFIED Code = iota
	RAW
	PROCESSING
	PROCESSING_ERROR
	READY
	INTERNAL_ERROR
)

type Result struct {
	data  []byte
	state Code
}

func NewResult(
	data []byte,
	state Code,
) *Result {
	return &Result{data, state}
}

var (
	NilResult   = NewResult(make([]byte, 0), UNSPECIFIED)
	EmptyResult = NewResult([]byte{}, RAW)
)

func (r Result) GetData() []byte {
	return r.data
}

func (r Result) GetState() int32 {
	return int32(r.state)
}

func (r *Result) SetData(data []byte) {
	r.data = data
}

func (r *Result) SetState(state Code) {
	r.state = state
}

func Equal(a, b *Result) bool {
	return bytes.Equal(a.GetData(), b.GetData()) && a.GetState() == b.GetState()
}
