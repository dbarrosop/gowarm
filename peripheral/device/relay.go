package device

import "machine"

type PinRelay struct {
	pin   machine.Pin
	state bool
}

func NewPinRelay(p machine.Pin) *PinRelay {
	p.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	return &PinRelay{p, false}
}

func (r *PinRelay) TurnOn() {
	r.state = true
	r.pin.High()
}

func (r *PinRelay) TurnOff() {
	r.state = false
	r.pin.Low()
}

func (r *PinRelay) On() bool {
	return r.state
}
