package terraform

type Root struct {
	Terraform Terraform `json:"terraform" hcle:"terraform,block"`
}

type Terraform struct {
	Backend Backend `json:"backend" hcle:"backend,block"`
}

type Backend interface {
	Name() string
}

type S3Backend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"s3" hcle:",body"`
}

func (s S3Backend) Name() string {
	return s.Type
}

type LocalBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"local" hcle:",body"`
}

func (s LocalBackend) Name() string {
	return s.Type
}

type RemoteBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"remote" hcle:",body"`
}

func (s RemoteBackend) Name() string {
	return s.Type
}

type AzurermBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"azurerm" hcle:",body"`
}

func (s AzurermBackend) Name() string {
	return s.Type
}

type ConsulBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"consul" hcle:",body"`
}

func (s ConsulBackend) Name() string {
	return s.Type
}

type CosBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"cos" hcle:",body"`
}

func (s CosBackend) Name() string {
	return s.Type
}

type GcsBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"gcs" hcle:",body"`
}

func (s GcsBackend) Name() string {
	return s.Type
}

type HttpBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"http" hcle:",body"`
}

func (s HttpBackend) Name() string {
	return s.Type
}

type KubernetesBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"kubernetes" hcle:",body"`
}

func (s KubernetesBackend) Name() string {
	return s.Type
}

type OSSBackend struct {
	Root
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"oss" hcle:",body"`
}

func (s OSSBackend) Name() string {
	return s.Type
}

type PGBackend struct {
	Type string         `json:"-" hcle:",label"`
	Data map[string]any `json:"pg" hcle:",body"`
}

func (s PGBackend) Name() string {
	return s.Type
}
func GetBackend(backendType string, config map[string]any) Root {

	var backend Backend
	switch backendType {
	case "local":
		backend = LocalBackend{
			Type: "local",
			Data: config,
		}
	case "remote":
		backend = RemoteBackend{
			Type: "remote",
			Data: config,
		}
	case "azurerm":
		backend = AzurermBackend{
			Type: "azurerm",
			Data: config,
		}
	case "consul":
		backend = ConsulBackend{
			Type: "consul",
			Data: config,
		}
	case "cos":
		backend = CosBackend{
			Type: "cos",
			Data: config,
		}
	case "gcs":
		backend = GcsBackend{
			Type: "gcs",
			Data: config,
		}
	case "http":
		backend = HttpBackend{
			Type: "http",
			Data: config,
		}
	case "kubernetes":
		backend = KubernetesBackend{
			Type: "kubernetes",
			Data: config,
		}
	case "oss":
		backend = OSSBackend{
			Type: "oss",
			Data: config,
		}
	case "pg":
		backend = PGBackend{
			Type: "pg",
			Data: config,
		}
	case "s3":
		backend = S3Backend{
			Type: "s3",
			Data: config,
		}
	}
	return Root{Terraform{Backend: backend}}
}
