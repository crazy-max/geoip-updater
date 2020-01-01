package maxmind

// Suffix represents the suffix of a database
type Suffix string

// Suffix enum
const (
	SfxTarGz = Suffix("tar.gz")
	SfxZip   = Suffix("zip")
)

// String returns the string representation of a suffix
func (sfx Suffix) String() string {
	return string(sfx)
}
