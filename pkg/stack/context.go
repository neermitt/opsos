package stack

import "context"

func SetStackName(ctx context.Context, stackName string) context.Context {
	return context.WithValue(ctx, "stackName", stackName)
}

func GetStackName(ctx context.Context) string {
	return ctx.Value("stackName").(string)
}

type Component struct {
	Type string
	Name string
}

func SetComponent(ctx context.Context, component Component) context.Context {
	return context.WithValue(ctx, "component", component)
}

func GetComponent(ctx context.Context) Component {
	return ctx.Value("component").(Component)
}

func SetComponentConfig(ctx context.Context, config *ConfigWithMetadata) context.Context {
	return context.WithValue(ctx, "component-config", config)
}

func GetComponentConfig(ctx context.Context) *ConfigWithMetadata {
	return ctx.Value("component-config").(*ConfigWithMetadata)
}
