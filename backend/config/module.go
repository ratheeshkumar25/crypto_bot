package config

import "go.uber.org/fx"

// Module provides the config module
var Module = fx.Options(
	fx.Provide(NewConfig),
)
