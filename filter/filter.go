package filter

type FilterInterface interface {
	Add(string)
	AddByte([]byte)
	Contains(string) bool
	ContainsByte([]byte) bool
}

func MakeFilter(length uint, name string) FilterInterface {
	switch name {
	case "local":
		return NewLocalFilter(length)
	case "redis":
		return NewRedisFilter()
	default:
		return NewLocalFilter(length)
	}
}
