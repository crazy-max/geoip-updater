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

// GetEditionID returns an edition ID from string
func GetEditionID(eidStr string) (EditionID, error) {
	eids := []EditionID{
		EIDGeoIP2City,
		EIDGeoIP2CityCSV,
		EIDGeoIP2Country,
		EIDGeoIP2CountryCSV,
		EIDGeoLite2ASN,
		EIDGeoLite2ASNCSV,
		EIDGeoLite2City,
		EIDGeoLite2CityCSV,
		EIDGeoLite2Country,
		EIDGeoLite2CountryCSV,
	}
	for _, eid := range eids {
		if EditionID(eidStr) == eid {
			return eid, nil
		}
	}
	return "", errors.Errorf("invalid edition ID: %s", eidStr)
}

// Suffix returns the suffix linked of an edition ID
func (eid EditionID) Suffix() Suffix {
	switch eid {
	case EIDGeoIP2City:
		return SfxTarGz
	case EIDGeoIP2CityCSV:
		return SfxZip
	case EIDGeoIP2Country:
		return SfxTarGz
	case EIDGeoIP2CountryCSV:
		return SfxZip
	case EIDGeoLite2ASN:
		return SfxTarGz
	case EIDGeoLite2ASNCSV:
		return SfxZip
	case EIDGeoLite2City:
		return SfxTarGz
	case EIDGeoLite2CityCSV:
		return SfxZip
	case EIDGeoLite2Country:
		return SfxTarGz
	case EIDGeoLite2CountryCSV:
		return SfxZip
	default:
		return ""
	}
}

// Filename returns the filename of an edition ID
func (eid EditionID) Filename() string {
	return fmt.Sprintf("%s.%s", eid.String(), eid.Suffix().String())
}

// String returns the string representation of an edition ID
func (eid EditionID) String() string {
	return string(eid)
}
