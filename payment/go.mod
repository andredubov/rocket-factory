module github.com/andredubov/rocket-factory/payment

go 1.24.4

replace github.com/andredubov/rocket-factory/shared => ../shared

require (
	github.com/andredubov/rocket-factory/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.73.0
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

require (
	github.com/andredubov/golibs v0.0.0-20240902121557-ded4e7068ebd
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
