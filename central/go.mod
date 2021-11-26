module github.com/dbarrosop/gowarm/central

go 1.16

require (
	github.com/brutella/hc v1.2.4
	github.com/dbarrosop/gowarm/peripheral v0.0.0-00010101000000-000000000000
	github.com/miekg/dns v1.1.41 // indirect
	github.com/pion/dtls/v2 v2.0.1-0.20200503085337-8e86b3a7d585
	github.com/plgd-dev/go-coap/v2 v2.4.0
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	tinygo.org/x/bluetooth v0.3.1-0.20210903130701-899467bab329
)

replace github.com/dbarrosop/gowarm/peripheral => ../peripheral
