package terraform

type Root struct {
	Terraform Terraform `json:"terraform" hcle:"terraform,block"`
}

type Terraform struct {
	Backend Backend `json:"-" hcle:"backend,block"`

	JSONBackend map[string]map[string]any `json:"backend,inline" `
}

type Backend struct {
	Type string         `hcle:",label"`
	Data map[string]any `hcle:",body"`
}

func GetBackend(backendType string, config map[string]any) Root {

	return Root{
		Terraform{
			Backend: Backend{
				Type: backendType,
				Data: config,
			},
			JSONBackend: map[string]map[string]any{backendType: config},
		}}
}
