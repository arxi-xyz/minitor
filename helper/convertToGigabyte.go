package helper

type Amount uint64

const (
	Byte Amount = 1

	KiloByte = 1024 * Byte
	MegaByte = 1024 * KiloByte
	GigaByte = 1024 * MegaByte
	TeraByte = 1024 * GigaByte
)

func Convert(from, to Amount, value float64) float64 {
	bytes := value * float64(from)
	return bytes / float64(to)
}
