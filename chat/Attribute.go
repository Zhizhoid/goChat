package main

type Attributes struct {
	attributes uint32
}

const (
	AttributeRead uint32 = 1 << iota
	AttributeWrite
	AttributeDelete
)

func (a *Attributes) SetAttributes(attrbutes uint32) {
	(*a).attributes |= attrbutes
}

func (a *Attributes) GetAttributes(attributes uint32) bool {
	return ((*a).attributes & attributes) > 0
}
