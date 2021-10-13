package thermostat

const (
	MODE_OFF  = 0x0
	MODE_HEAT = 0x1
	MODE_COOL = 0x2
	MODE_AUTO = 0x3
)

type Sensor interface {
	// Temperature returns the temperature in celsius
	Temperature() (float32, error)
	// Humidity returns the relative humidity
	Humidity() (float32, error)
}

type Relay interface {
	// Turn on relay
	TurnOn()
	// Turn off relay
	TurnOff()
	// Returns whether it's on or not
	On() bool
}

type Comms interface {
	SendTemperature(float32) error
	SendHumidity(float32) error
}

type Thermostat struct {
	sensor           Sensor
	relay            Relay
	targetTemp       float32
	hysteresisMargin float32
	mode             byte
}

func New(sensor Sensor, relay Relay, targetTemp, hysteresisMargin float32) *Thermostat {
	return &Thermostat{
		sensor:           sensor,
		relay:            relay,
		targetTemp:       targetTemp,
		hysteresisMargin: hysteresisMargin,
	}
}

func (th *Thermostat) Process() (float32, float32) {
	temp, err := th.sensor.Temperature()
	if err != nil {
		println("got error reading temperature: %s", err.Error())
		return 0, 0
	}

	humidity, err := th.sensor.Humidity()
	if err != nil {
		println("got error reading humidity: %s", err.Error())
		return 0, 0
	}

	thermostatToogle(temp, th.targetTemp, th.hysteresisMargin, th.relay, th.mode)

	return temp, humidity
}

func thermostatToogle(currentTemp, targetTemp, hysteresisMargin float32, relay Relay, mode byte) {
	switch {
	case mode == MODE_OFF:
		relay.TurnOff()
	case !relay.On() && currentTemp <= targetTemp-hysteresisMargin:
		relay.TurnOn()
	case relay.On() && currentTemp >= targetTemp+hysteresisMargin:
		relay.TurnOff()
	}
}

func (th *Thermostat) SetTargetTemperature(value float32) {
	th.targetTemp = value
}

func (th *Thermostat) SetMode(mode byte) {
	th.mode = mode
}
