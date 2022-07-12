package tradfri

import "time"

type GatewayInfo struct {
	ID string `json:"9081"`

	Firmware string `json:"9029"`

	Time       time.Time `json:"9060"`
	TimeServer string    `json:"9023"`

	EpochNow     int64 `json:"9059"`
	EpochCreated int64 `json:"9069"`
}

func (g *GatewayInfo) Now() time.Time {
	return time.Unix(g.EpochNow, 0)
}

func (g *GatewayInfo) CreatedAt() time.Time {
	return time.Unix(g.EpochCreated, 0)
}

type DeviceInfo struct {
	ID   int    `json:"9003"` // Device ID
	Name string `json:"9001"` // Device Name

	Type DeviceType `json:"5750"` // Application Type

	Metadata struct {
		Manufacturer string      `json:"0"`
		Model        string      `json:"1"`
		Serial       string      `json:"2"`
		Firmware     string      `json:"3"`
		PowerSource  PowerSource `json:"6"`
		Battery      int         `json:"9"`
	} `json:"3"`

	Speaker *struct {
		Devices []string `json:"9115"`
	} `json:"15017"`

	LightSettings   []LightSettings   `json:"3311"`
	OutletSettings  []OutletSettings  `json:"3312"`
	SwitchSettings  []SwitchSettings  `json:"15009"`
	SpeakerSettings []SpeakerSettings `json:"15018"`
	Sensors         []Sensor          `json:"3300"`

	EpochCreated   int64 `json:"9002"`
	EpochUpdated   int64 `json:"9020"` // Last Seen?
	Alive          int   `json:"9019"`
	OtaUpdateState int   `json:"9054"`
}

type Sensor struct {
	AppType                 string  `json:"5750"`
	SensorType              string  `json:"5751"`
	MinMeasuredValue        int64   `json:"5601"`
	MaxMeasuredValue        int64   `json:"5602"`
	MinRangeValue           int64   `json:"5603"`
	MaxRangeValue           int64   `json:"5604"`
	ResetMinMaxMeasureValue bool    `json:"5605"`
	SensorValue             float64 `json:"5700"`
	Unit                    string  `json:"5701"`
}

type DeviceType int

const (
	DeviceTypeSwitch         DeviceType = 0 // Normal remotes
	DeviceTypeSwitchSlave               = 1 // A remote paired with another remote (https://www.reddit.com/r/tradfri/comments/6x1miq)
	DeviceTypeBulb                      = 2
	DeviceTypeControlOutlet             = 3
	DeviceTypeSensor                    = 4
	DeviceTypeSignalRepeater            = 6
	DeviceTypeBlind                     = 7
	DeviceTypeSoundRemote               = 8 // Symfonisk
	DeviceTypeAirPurifier               = 9 // Starkvind
)

type PowerSource int

const (
	PowerSourceDC              PowerSource = 0
	PowerSourceBatteryInternal             = 1
	PowerSourceBatteryExternal             = 2
	PowerSourceBattery                     = 3
	PowerSourceEthernet                    = 4
	PowerSourceUSB                         = 5
	PowerSourceAC                          = 6
	PowerSourceSolar                       = 7
)

type DeviceSettings struct {
	LightSettings  []LightSettings  `json:"3311,omitempty"`
	OutletSettings []OutletSettings `json:"3312,omitempty"`
}

type LightSettings struct {
	Power  *int `json:"5850,omitempty"`
	Dimmer *int `json:"5851,omitempty"`

	ColorHue        *int    `json:"5707,omitempty"`
	ColorSaturation *int    `json:"5708,omitempty"`
	ColorX          *int    `json:"5709,omitempty"`
	ColorY          *int    `json:"5710,omitempty"`
	Color           *string `json:"5706,omitempty"`

	//ColorTemperature *int    `json:"5711,omitempty"`
	Duration *int `json:"5712,omitempty"`

	Device *int `json:"9003,omitempty"`
}

type OutletSettings struct {
	Power  *int `json:"5850,omitempty"`
	Dimmer *int `json:"5851,omitempty"`

	Device *int `json:"9003,omitempty"`
}

type SwitchSettings struct {
	Device *int `json:"9003,omitempty"`
}

type SpeakerSettings struct {
	Device *int `json:"9003,omitempty"`
}

func (d *DeviceInfo) CreatedAt() time.Time {
	return time.Unix(d.EpochCreated, 0)
}

func (d *DeviceInfo) UpdatedAt() time.Time {
	return time.Unix(d.EpochUpdated, 0)
}

type GroupInfo struct {
	ID   int    `json:"9003"` // Group ID
	Name string `json:"9001"` // Group Name

	EpochCreated int64 `json:"9002"`
}

func (g *GroupInfo) CreatedAt() time.Time {
	return time.Unix(g.EpochCreated, 0)
}

type SceneInfo struct {
	ID   int    `json:"9003"` // Group ID
	Name string `json:"9001"` // Group Name

	EpochCreated int64 `json:"9002"`

	LightSettings []LightSettings `json:"15013"`
}

func (s *SceneInfo) CreatedAt() time.Time {
	return time.Unix(s.EpochCreated, 0)
}
