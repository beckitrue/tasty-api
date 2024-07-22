module tastyapi

go 1.21.4

require (
	github.com/beckitrue/tasty-api v0.0.0-20240630181236-8c871ae27239
	github.com/beckitrue/tastyapi/cmd/tastymenu v0.0.0-20240630181236-8c871ae27239
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.27.2 // indirect
	github.com/xrash/smetrics v0.0.0-20240312152122-5f08fbb34913 // indirect
)

replace github.com/beckitrue/tastyapi/cmd/tastymenu => ./cmd/tastymenu
