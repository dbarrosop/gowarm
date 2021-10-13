package thermostat

import "testing"

type dummyRelay struct {
	on      bool
	toggled bool
}

func (r *dummyRelay) TurnOn() {
	r.toggled = !r.on
	r.on = true
}

func (r *dummyRelay) TurnOff() {
	r.toggled = r.on
	r.on = false
}

func (r *dummyRelay) On() bool {
	return r.on
}

func TestThermostatToggle(t *testing.T) {
	cases := []struct {
		name               string
		currentTemp        float32
		targetTemp         float32
		hysteresisMarging  float32
		initialRelayState  bool
		expectedRelayState bool
		expectedToggle     bool
	}{
		{
			name:               "20.0, 0.2, 22.0, on",
			currentTemp:        20.0,
			targetTemp:         22.0,
			hysteresisMarging:  0.2,
			initialRelayState:  true,
			expectedRelayState: true,
			expectedToggle:     false,
		},
		{
			name:               "20.0, 0.2, 22.0, off",
			currentTemp:        20.0,
			targetTemp:         22.0,
			hysteresisMarging:  0.1,
			initialRelayState:  false,
			expectedRelayState: true,
			expectedToggle:     true,
		},
		{
			name:               "21.9, 0.2, 22.0, off",
			currentTemp:        21.9,
			targetTemp:         22.0,
			hysteresisMarging:  0.2,
			initialRelayState:  false,
			expectedRelayState: false,
			expectedToggle:     false,
		},
		{
			name:               "22.5, 0.2, 22.0, off",
			currentTemp:        22.5,
			targetTemp:         22.0,
			hysteresisMarging:  0.2,
			initialRelayState:  false,
			expectedRelayState: false,
			expectedToggle:     false,
		},
		{
			name:               "22.5, 0.2, 22.0, on",
			currentTemp:        22.5,
			targetTemp:         22.0,
			hysteresisMarging:  0.2,
			initialRelayState:  true,
			expectedRelayState: false,
			expectedToggle:     true,
		},
		{
			name:               "22.1, 0.2, 22.0, off",
			currentTemp:        22.1,
			targetTemp:         22.0,
			hysteresisMarging:  0.2,
			initialRelayState:  false,
			expectedRelayState: false,
			expectedToggle:     false,
		},
		{
			name:               "22.1, 0.2, 22.0, on",
			currentTemp:        22.1,
			targetTemp:         22.0,
			hysteresisMarging:  0.2,
			initialRelayState:  true,
			expectedRelayState: true,
			expectedToggle:     false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc := tc

			relay := &dummyRelay{tc.initialRelayState, false}
			thermostatToogle(tc.currentTemp, tc.targetTemp, tc.hysteresisMarging, relay, 1)

			if relay.On() != tc.expectedRelayState {
				t.Errorf("Current state (%t) != expected state (%t)", relay.On(), tc.expectedRelayState)
			}
			if relay.toggled != tc.expectedToggle {
				t.Errorf("toggle (%t) != expected toggle (%t)", relay.toggled, tc.expectedToggle)
			}
		})
	}
}
