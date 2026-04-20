package maxmind

import (
	"fmt"

	"github.com/pkg/errors"
)

// EditionID represents the edition ID of a database
type EditionID string

// Edition ID enum
const (
	EIDGeoIP2City         = EditionID("GeoIP2-City")
	EIDGeoIP2CityCSV      = EditionID("GeoIP2-City-CSV")
	EIDGeoIP2Country      = EditionID("GeoIP2-Country")
	EIDGeoIP2CountryCSV   = EditionID("GeoIP2-Country-CSV")
	EIDGeoLite2ASN        = EditionID("GeoLite2-ASN")
	EIDGeoLite2ASNCSV     = EditionID("GeoLite2-ASN-CSV")
	EIDGeoLite2City       = EditionID("GeoLite2-City")
	EIDGeoLite2CityCSV    = EditionID("GeoLite2-City-CSV")
	EIDGeoLite2Country    = EditionID("GeoLite2-Country")
	EIDGeoLite2CountryCSV = EditionID("GeoLite2-Country-CSV")
)

// Suffix represents the suffix of a database
type Suffix string

// Suffix enum
const (
	SfxTarGz = Suffix("tar.gz")
	SfxZip   = Suffix("zip")
)

var editionIDSuffixes = map[EditionID]Suffix{
	EIDGeoIP2City:         SfxTarGz,
	EIDGeoIP2CityCSV:      SfxZip,
	EIDGeoIP2Country:      SfxTarGz,
	EIDGeoIP2CountryCSV:   SfxZip,
	EIDGeoLite2ASN:        SfxTarGz,
	EIDGeoLite2ASNCSV:     SfxZip,
	EIDGeoLite2City:       SfxTarGz,
	EIDGeoLite2CityCSV:    SfxZip,
	EIDGeoLite2Country:    SfxTarGz,
	EIDGeoLite2CountryCSV: SfxZip,
}

// GetEditionID returns an edition ID from string
func GetEditionID(eidStr string) (EditionID, error) {
	eid := EditionID(eidStr)
	if _, ok := editionIDSuffixes[eid]; ok {
		return eid, nil
	}
	return "", errors.Errorf("invalid edition ID: %s", eidStr)
}

// Suffix returns the suffix linked of an edition ID
func (eid EditionID) Suffix() Suffix {
	return editionIDSuffixes[eid]
}

// Filename returns the filename of an edition ID
func (eid EditionID) Filename() string {
	return fmt.Sprintf("%s.%s", eid.String(), eid.Suffix().String())
}

// String returns the string representation of an edition ID
func (eid EditionID) String() string {
	return string(eid)
}

// String returns the string representation of a suffix
func (sfx Suffix) String() string {
	return string(sfx)
}
