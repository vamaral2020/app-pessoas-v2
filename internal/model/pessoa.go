package model

import "time"

type Pessoa struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Email        string    `json:"email"`
	DataCadastro time.Time `json:"data_cadastro"`
}

type PessoaCreateRequest struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
}
