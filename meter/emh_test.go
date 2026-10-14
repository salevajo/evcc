package meter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEMH(t *testing.T) {
	body, err := os.ReadFile("emh/start.extended.json")
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != "u" || p != "p" || r.URL.Path != "/json/metering/origin/1efr0123456789/extended" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	m, err := NewEMH(srv.URL, "u", "p", "1efr0123456789", time.Second)
	require.NoError(t, err)

	p, err := m.CurrentPower()
	require.NoError(t, err)
	assert.Equal(t, 1392.0, p)

	e, err := m.TotalEnergy()
	require.NoError(t, err)
	assert.InDelta(t, 1918.6528, e, 1e-6)

	l1, _, _, err := m.Voltages()
	require.NoError(t, err)
	assert.InDelta(t, 231.8, l1, 1e-6)

	_, _, l3, err := m.Currents()
	require.NoError(t, err)
	assert.InDelta(t, 3.48, l3, 1e-6)
}
