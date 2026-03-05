package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"app-pessoas/internal/model"
	"app-pessoas/internal/repo"
)

type Handler struct {
	Pessoas *repo.PessoaRepo
}

func NewHandler(pessoas *repo.PessoaRepo) *Handler {
	return &Handler{Pessoas: pessoas}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /pessoas", h.createPessoa)
	mux.HandleFunc("GET /pessoas/{id}", h.getPessoaByID)
	mux.HandleFunc("GET /pessoas", h.listPessoas)

}

func (h *Handler) createPessoa(writer http.ResponseWriter, request *http.Request) {
	var pessoaRequest model.PessoaCreateRequest
	if err := json.NewDecoder(request.Body).Decode(&pessoaRequest); err != nil {
		fmt.Println(err)
		writeJSON(writer, http.StatusBadRequest, map[string]any{"error": "Json invalido"})
		return
	}
	pessoaRequest.Nome = strings.TrimSpace(pessoaRequest.Nome)
	pessoaRequest.Email = strings.TrimSpace(pessoaRequest.Email)

	if pessoaRequest.Nome == "" || pessoaRequest.Email == "" {
		writeJSON(writer, http.StatusBadRequest, map[string]any{
			"error": "Nome e Email são obrigatórios",
		})
		return
	}

	if !strings.Contains(pessoaRequest.Email, "@") {
		writeJSON(writer, http.StatusBadRequest, map[string]any{
			"error": "Email invalido",
		})
		return
	}
	ctx, cancel := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancel()

	p, err := h.Pessoas.Create(ctx, pessoaRequest)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, repo.ErrEmailDuplicado) {
			writeJSON(writer, http.StatusConflict, map[string]any{
				"error": "Email duplicado",
			})
			return
		}
		writeJSON(writer, http.StatusInternalServerError, map[string]any{"error": "Falha ao salvar pessoa no banco de dados"})
		return
	}
	writeJSON(writer, http.StatusCreated, p)
}

func writeJSON(writer http.ResponseWriter, status int, v any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(v)
}

func (h *Handler) getPessoaByID(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		writeJSON(writer, http.StatusBadRequest, map[string]any{"error": "Id invalido"})
		return
	}
	ctx, cancel := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancel()
	p, err := h.Pessoas.GetByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, repo.ErrNaoEncontrado) {
			writeJSON(writer, http.StatusNotFound, map[string]any{"error": "Pessoa nao encontrado"})
			return
		}
		writeJSON(writer, http.StatusInternalServerError, map[string]any{"error": "Erro ao consultar pessoa"})
	}
	writeJSON(writer, http.StatusOK, p)
}

func (h *Handler) listPessoas(writer http.ResponseWriter, request *http.Request) {
	limit := 50
	offset := 0

	if v := request.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n < 200 {
			limit = n
		}
	}
	if v := request.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			offset = n
		}
	}
	ctx, cancel := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancel()

	pessoas, err := h.Pessoas.List(ctx, limit, offset)
	if err != nil {
		fmt.Println("Erro ao listar pessoas", err)
		writeJSON(writer, http.StatusInternalServerError, map[string]any{
			"error": "Erro ao listar pessoas",
		})
		return
	}

	writeJSON(writer, http.StatusOK, map[string]any{
		"pessoas": pessoas,
		"limit":   limit,
		"offset":  offset,
	})
}
