module go_agent

go 1.22.5

require (
	github.com/bluenviron/goroslib/v2 v2.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/gookit/color v1.5.4 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

replace github.com/bluenviron/goroslib/v2 => github.com/SevenRobotics/goroslib/v2 v2.0.0-20240802101958-e71c19414184
