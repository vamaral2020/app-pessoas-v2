package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"

	"app-pessoas/internal/model"
)

type PessoaRepo struct {
	DB *sql.DB
}

func NewPessoaRepo(db *sql.DB) *PessoaRepo {
	return &PessoaRepo{DB: db}
}

func (r *PessoaRepo) Create(ctx context.Context, req model.PessoaCreateRequest) (model.Pessoa, error) {

	var p model.Pessoa

	err := r.DB.QueryRowContext(ctx, `
    INSERT INTO pessoas (nome, email)
    VALUES ($1, $2)
    RETURNING id, nome, email, criado_em
`, req.Nome, req.Email).Scan(&p.ID, &p.Nome, &p.Email, &p.DataCadastro)

	if err != nil {
		fmt.Println(err)
		if isUniqueViolation(err) {
			return model.Pessoa{}, ErrEmailDuplicado
		}
		return model.Pessoa{}, err
	}
	return p, nil
}

func (r *PessoaRepo) GetByID(ctx context.Context, id int64) (model.Pessoa, error) {
	var p model.Pessoa
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, nome, email, criado_em
		FROM pessoas
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Nome, &p.Email, &p.DataCadastro)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Pessoa{}, ErrNaoEncontrado
		}
		return model.Pessoa{}, err
	}

	return p, nil
}

var (
	ErrNaoEncontrado  = errors.New("nao encontrado")
	ErrEmailDuplicado = errors.New("email duplicatado")
)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *PessoaRepo) List(ctx context.Context, limit, offset int) ([]model.Pessoa, error) {

	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, nome, email, criado_em
		FROM pessoas
		ORDER BY id
		LIMIT $1 OFFSET $2
		`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pessoas []model.Pessoa
	for rows.Next() {
		var p model.Pessoa
		if err := rows.Scan(&p.ID, &p.Nome, &p.Email, &p.DataCadastro); err != nil {
			return nil, err
		}
		pessoas = append(pessoas, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pessoas, nil
}
