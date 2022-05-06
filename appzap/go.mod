module github.com/yuridevx/app/appzap

go 1.18

replace github.com/yuridevx/app => ../

require (
	github.com/yuridevx/app v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.21.0
)

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)