package meter

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
	"github.com/evcc-io/evcc/util/transport"
)

// EMH meter implementation for EMH smart meter gateways (CASA) using the local JSON interface
type EMH struct {
	dataG func() (map[string]float64, error)
}

type emhResponse struct {
	Values []struct {
		LogicalName string `json:"logical_name"`
		Scaler      int
		Value       string
	}
}

// OBIS registers (class 1, media 0) as delivered by the gateway
const (
	emhEnergyImport = "0100010800ff"
	emhPower        = "0100100700ff"
)

var (
	emhCurrents = [3]string{"01001f0700ff", "0100330700ff", "0100470700ff"}
	emhVoltages = [3]string{"0100200700ff", "0100340700ff", "0100480700ff"}
	emhPowers   = [3]string{"0100240700ff", "0100380700ff", "01004c0700ff"}
)

func init() {
	registry.Add("emh", NewEMHFromConfig)
}

// NewEMHFromConfig creates an EMH meter from generic config
func NewEMHFromConfig(other map[string]any) (api.Meter, error) {
	cc := struct {
		URI      string
		User     string
		Password string
		ID       string
		Cache    time.Duration
	}{
		Cache: 30 * time.Second,
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	return NewEMH(cc.URI, cc.User, cc.Password, cc.ID, cc.Cache)
}

// NewEMH creates an EMH smart meter gateway meter
func NewEMH(uri, user, password, id string, cache time.Duration) (*EMH, error) {
	if uri == "" || id == "" {
		return nil, errors.New("missing uri or id")
	}

	log := util.NewLogger("emh")
	client := request.NewHelper(log)
	client.Client.Transport = request.NewTripper(log,
		transport.BasicAuth(user, password, transport.Insecure()))

	url := fmt.Sprintf("%s/json/metering/origin/%s/extended",
		util.DefaultScheme(strings.TrimRight(uri, "/"), "https"), id)

	return &EMH{
		dataG: util.Cached(func() (map[string]float64, error) {
			var res emhResponse
			if err := client.GetJSON(url, &res); err != nil {
				return nil, err
			}

			values := make(map[string]float64, len(res.Values))
			for _, v := range res.Values {
				raw, err := strconv.ParseFloat(v.Value, 64)
				if err != nil {
					continue
				}
				// logical name is "<register>.<meter id>.sm"
				register, _, _ := strings.Cut(v.LogicalName, ".")
				values[register] = raw * math.Pow10(v.Scaler)
			}
			return values, nil
		}, cache),
	}, nil
}

func (m *EMH) get(register string) (float64, error) {
	data, err := m.dataG()
	if err != nil {
		return 0, err
	}
	v, ok := data[register]
	if !ok {
		return 0, api.ErrNotAvailable
	}
	return v, nil
}

func (m *EMH) getPhases(registers [3]string) (float64, float64, float64, error) {
	var res [3]float64
	for i, r := range registers {
		v, err := m.get(r)
		if err != nil {
			return 0, 0, 0, err
		}
		res[i] = v
	}
	return res[0], res[1], res[2], nil
}

var _ api.Meter = (*EMH)(nil)

// CurrentPower implements the api.Meter interface
func (m *EMH) CurrentPower() (float64, error) {
	return m.get(emhPower)
}

var _ api.MeterEnergy = (*EMH)(nil)

// TotalEnergy implements the api.MeterEnergy interface (register is delivered in Wh)
func (m *EMH) TotalEnergy() (float64, error) {
	v, err := m.get(emhEnergyImport)
	return v / 1e3, err
}

var _ api.PhaseCurrents = (*EMH)(nil)

// Currents implements the api.PhaseCurrents interface
func (m *EMH) Currents() (float64, float64, float64, error) {
	return m.getPhases(emhCurrents)
}

var _ api.PhaseVoltages = (*EMH)(nil)

// Voltages implements the api.PhaseVoltages interface
func (m *EMH) Voltages() (float64, float64, float64, error) {
	return m.getPhases(emhVoltages)
}

var _ api.PhasePowers = (*EMH)(nil)

// Powers implements the api.PhasePowers interface
func (m *EMH) Powers() (float64, float64, float64, error) {
	return m.getPhases(emhPowers)
}
