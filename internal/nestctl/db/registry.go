package db

type RegistryOption []byte

type Registry struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	CaCert   string `json:"ca-cert"`
	Registry string `json:"registry"`
}
