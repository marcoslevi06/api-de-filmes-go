# API de Filmes

CRUD de filmes em Go, persistido em MongoDB, com arquitetura hexagonal
(ports & adapters).

## Como rodar

Pré-requisito: Docker instalado e com o daemon rodando (não precisa de Go
nem MongoDB instalados na máquina).

```bash
git clone <url-do-repositório>
cd sipub_teste
./start.sh
```

O script sobe três containers — API (porta `8080`), MongoDB (porta
`27017`) e a documentação Swagger (porta `8081`) — carrega os filmes
iniciais no banco e deixa a API pronta em `http://localhost:8080`.
`Ctrl+C` encerra e remove os três containers.

Os dados não persistem entre execuções: a cada `./start.sh`, o Mongo sobe
vazio e é populado de novo a partir de `movies.json`.

## Estrutura do projeto

```
sipub_teste/
├── start.sh                 # sobe tudo com um comando
├── docker/
│   ├── api/Dockerfile
│   └── mongo/Dockerfile
├── docs/
│   ├── fluxograma.drawio     # diagrama da arquitetura e do fluxo
│   ├── testes-curl.md        # roteiro de testes manuais via curl
│   └── swagger.yaml          # especificação OpenAPI da API
└── api/
    ├── cmd/api/main.go
    ├── movies.json
    └── internal/
        ├── domain/           # entidade Movie + interface MovieRepository
        ├── usecase/          # MovieService (regra de negócio)
        └── adapter/
            ├── http/          # recebe as requisições HTTP
            ├── mongo/         # persiste no MongoDB
            └── memory/        # versão em memória, sem banco
```

## Arquitetura

O projeto segue arquitetura hexagonal: a regra de negócio
(`MovieService`) não conhece HTTP nem MongoDB, só a interface
`MovieRepository`. Quem implementa essa interface (Mongo, memória) é um
detalhe que pode ser trocado sem mexer na regra de negócio nem no handler
HTTP.

Diagrama completo em [`docs/fluxograma.drawio`](docs/fluxograma.drawio)
(abra em [app.diagrams.net](https://app.diagrams.net) ou na extensão
"Draw.io Integration" do VS Code).

## Rotas

| Método | Rota | Descrição |
|---|---|---|
| GET | `/movies` | Lista todos os filmes |
| GET | `/movies/{id}` | Busca um filme pelo id |
| POST | `/movies` | Cria um filme |
| PUT | `/movies/{id}` | Atualiza um filme (extra, não pedido no enunciado) |
| DELETE | `/movies/{id}` | Remove um filme |

Exemplos de uso com curl (requisição, resposta e erros) em
[`docs/testes-curl.md`](docs/testes-curl.md).

Especificação OpenAPI/Swagger completa (schemas de request/response) em
[`docs/swagger.yaml`](docs/swagger.yaml). O `./start.sh` já sobe a
documentação interativa em `http://localhost:8081` automaticamente —
não precisa de nenhum passo extra.

Alternativa sem Docker: colar o conteúdo do arquivo em
[editor.swagger.io](https://editor.swagger.io).

## Variáveis de ambiente - (Já definidas como padrão no arquivo start.sh)
| Variável | Padrão |
|---|---|
| `MONGO_DB` | `sipub` |
| `MONGO_COLLECTION` | `movies` |
| `MONGO_URI` | `mongodb://sipub-mongo:27017` |
