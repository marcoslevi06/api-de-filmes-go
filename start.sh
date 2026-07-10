#!/usr/bin/env bash
# Builda a imagem do MongoDB, sobe um container novo e efêmero do banco
# (sem volume — os dados duram só enquanto o script estiver rodando) e, em
# seguida, executa a API Go apontando para ele. Ao encerrar (Ctrl+C), o
# container do Mongo é removido automaticamente. Uso: ./build_projeto.sh

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DIR="$ROOT_DIR/api"
MONGO_DOCKERFILE="$ROOT_DIR/docker/mongo/Dockerfile"
MONGO_IMAGE="sipub-mongo:local"
MONGO_CONTAINER="sipub-mongo"
MONGO_PORT="27017"


export MONGO_URI="${MONGO_URI:-mongodb://localhost:${MONGO_PORT}}"
export MONGO_DB="${MONGO_DB:-sipub}"
export MONGO_COLLECTION="${MONGO_COLLECTION:-movies}"

cleanup() {
  echo
  echo "==> Encerrando e removendo o container do MongoDB..."
  docker rm -f "$MONGO_CONTAINER" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "==> Construindo imagem do MongoDB..."
docker build -t "$MONGO_IMAGE" -f "$MONGO_DOCKERFILE" "$ROOT_DIR"

echo "==> Subindo container novo do MongoDB..."
docker rm -f "$MONGO_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$MONGO_CONTAINER" -p "${MONGO_PORT}:27017" "$MONGO_IMAGE" >/dev/null
echo "    Container '$MONGO_CONTAINER' criado."

echo "==> Aguardando o MongoDB responder..."
for i in $(seq 1 30); do
  if docker exec "$MONGO_CONTAINER" mongosh --quiet --eval 'db.runCommand({ ping: 1 })' >/dev/null 2>&1; then
    echo "    MongoDB pronto."
    break
  fi
  if [ "$i" -eq 30 ]; then
    echo "Erro: MongoDB não respondeu a tempo." >&2
    exit 1
  fi
  sleep 1
done

echo "==> Subindo a API Go em http://localhost:8080 (Ctrl+C para encerrar)..."
cd "$API_DIR"
go run ./cmd/api
