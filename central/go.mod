module github.com/dbarrosop/gowarm/central

go 1.16

require (
	github.com/dbarrosop/gowarm/peripheral v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.8.1
	tinygo.org/x/bluetooth v0.3.0
)

replace github.com/dbarrosop/gowarm/peripheral => ../peripheral
