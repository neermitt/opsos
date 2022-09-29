package config

import "context"

func SetConfig(ctx context.Context, conf *Configuration) context.Context {
	return context.WithValue(ctx, "config", conf)
}

func GetConfig(ctx context.Context) *Configuration {
	return ctx.Value("config").(*Configuration)
}
