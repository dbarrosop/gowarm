module github.com/dbarrosop/gowarm/central

go 1.16

require (
	github.com/brutella/hc v1.2.4
	github.com/dbarrosop/gowarm/peripheral v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	tinygo.org/x/bluetooth v0.3.0
)

replace github.com/dbarrosop/gowarm/peripheral => ../peripheral
