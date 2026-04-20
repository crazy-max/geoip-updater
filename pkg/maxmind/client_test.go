package maxmind

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	cases := []struct {
		name     string
		wantData Config
		wantErr  bool
	}{
		{
			name: "Empty license key",
			wantData: Config{
				LicenseKey: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid base URL",
			wantData: Config{
				LicenseKey: "0123456789",
				BaseURL:    "foo.bar",
			},
			wantErr: true,
		},
		{
			name: "Success",
			wantData: Config{
				LicenseKey: "0123456789",
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(context.Background(), tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestNewUsesProvidedContextAndHTTPClient(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	httpClient := &http.Client{}
	client, err := New(ctx, Config{
		HTTPClient: httpClient,
		LicenseKey: "0123456789",
	})
	require.NoError(t, err)

	assert.Same(t, ctx, client.ctx)
	assert.Same(t, httpClient, client.http)
}

func TestExpectedHashUsesConfiguredContext(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("should not be reached"))
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(context.Canceled)

	client, err := New(ctx, Config{
		HTTPClient: srv.Client(),
		LicenseKey: "0123456789",
		BaseURL:    srv.URL,
	})
	require.NoError(t, err)

	downloader, err := client.NewDownloader(EIDGeoLite2ASNCSV, t.TempDir())
	require.NoError(t, err)

	_, err = downloader.expectedHash()
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}
