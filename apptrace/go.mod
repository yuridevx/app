module github.com/yuridevx/app/apptrace

go 1.18

replace github.com/yuridevx/app => ../

require (
	github.com/newrelic/go-agent/v3 v3.15.2
	github.com/yuridevx/app v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.21.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac // indirect
	google.golang.org/grpc v1.45.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
