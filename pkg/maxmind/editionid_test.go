package maxmind

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEditionID(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    EditionID
		wantSfx Suffix
		wantErr bool
	}{
		{
			name:    "Valid MMDB edition",
			input:   "GeoLite2-City",
			want:    EIDGeoLite2City,
			wantSfx: SfxTarGz,
		},
		{
			name:    "Valid CSV edition",
			input:   "GeoLite2-City-CSV",
			want:    EIDGeoLite2CityCSV,
			wantSfx: SfxZip,
		},
		{
			name:    "Invalid edition",
			input:   "GeoLite2-Unknown",
			wantErr: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEditionID(tt.input)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantSfx, got.Suffix())
		})
	}
}
