package mcp

type AccessRoute struct {
	Sts  station
	Code Code
}

func (r *AccessRoute) BinaryRoute() []byte {

	return nil
}

func (r *AccessRoute) AsciiRoute() []byte {

	return nil
}

func (r *AccessRoute) Len() int64 {
	return 0
}
