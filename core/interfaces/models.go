package interfaces

type Transformation interface {
	GetUrl() string
	GetWidth() int32
	GetHeight() int32
}

type Result interface {
	GetData() []byte
	GetState() int32
}
