package dto

type HashRequest struct {
	SecretName string `json:"secretName"`
}

type HashResponse struct {
	SecretNameHash string `json:"secretNameHash"`
}