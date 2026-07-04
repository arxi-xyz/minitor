package collector

type Metric struct {
	Type  string
	Value int

	AllResource  int
	UsedResource int
}
