module github.com/andredubov/rocket-factory/payment

go 1.24.4

replace github.com/andredubov/rocket-factory/shared => ../shared

replace github.com/andredubov/rocket-factory/platform => ../platform

require (
	github.com/andredubov/rocket-factory/platform v0.0.0-00010101000000-000000000000
	github.com/andredubov/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.74.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/caarlos0/env/v11 v11.3.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/sys v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
