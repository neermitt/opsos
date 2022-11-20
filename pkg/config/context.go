package config

import (
	"context"

	v1 "github.com/neermitt/opsos/api/v1"
)

func SetConfig(ctx context.Context, conf *v1.ConfigSpec) context.Context {
	return context.WithValue(ctx, "config", conf)
}

func GetConfig(ctx context.Context) *v1.ConfigSpec {
	return ctx.Value("config").(*v1.ConfigSpec)
}
