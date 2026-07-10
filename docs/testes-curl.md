# Roteiro de testes manuais (curl)

Pré-requisito: a API rodando (`./start.sh`), disponível em
`http://localhost:8080`.

## 1. Listar todos os filmes

```bash
curl -s http://localhost:8080/movies | head -c 300
```

Esperado: `200 OK`, lista JSON com os filmes carregados de `movies.json`.

## 2. Buscar um filme existente

```bash
curl -i http://localhost:8080/movies/8
```

Esperado: `200 OK`

```json
{"id":8,"title":"Edison Kinetoscopic Record of a Sneeze (1894)","year":"1894"}
```

## 3. Buscar um filme inexistente

```bash
curl -i http://localhost:8080/movies/999999999
```

Esperado: `404 Not Found` — "Filme não encontrado - Tente buscar outro ID."

## 4. Buscar com id inválido

```bash
curl -i http://localhost:8080/movies/abc
```

Esperado: `400 Bad Request` — "ID inválido."

## 5. Criar um filme

```bash
curl -i -X POST http://localhost:8080/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "Duna", "year": "2021"}'
```

Esperado: `201 Created`

```json
{"id": <novo-id>, "title": "Duna", "year": "2021"}
```

Guarde o `id` retornado — ele é usado nos próximos passos.

## 6. Criar com corpo inválido

```bash
curl -i -X POST http://localhost:8080/movies \
  -H "Content-Type: application/json" \
  -d 'isso não é um json'
```

Esperado: `400 Bad Request` — "Estrutura inválida."

## 7. Atualizar o filme criado

```bash
curl -i -X PUT http://localhost:8080/movies/<id-do-passo-5> \
  -H "Content-Type: application/json" \
  -d '{"title": "Duna: Parte 2", "year": "2024"}'
```

Esperado: `200 OK`

```json
{"id": <id>, "title": "Duna: Parte 2", "year": "2024"}
```

## 8. Atualizar um filme inexistente

```bash
curl -i -X PUT http://localhost:8080/movies/999999999 \
  -H "Content-Type: application/json" \
  -d '{"title": "X", "year": "2024"}'
```

Esperado: `404 Not Found`

## 9. Remover o filme criado

```bash
curl -i -X DELETE http://localhost:8080/movies/<id-do-passo-5>
```

Esperado: `204 No Content`, sem corpo na resposta.

## 10. Confirmar a remoção

```bash
curl -i http://localhost:8080/movies/<id-do-passo-5>
```

Esperado: `404 Not Found` — o filme não existe mais.

## 11. Remover um filme inexistente

```bash
curl -i -X DELETE http://localhost:8080/movies/999999999
```

Esperado: `404 Not Found`
