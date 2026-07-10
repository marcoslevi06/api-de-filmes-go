#!/usr/bin/env bash
# Builda as imagens do MongoDB e da API, sobe os dois em containers
# efêmeros na mesma rede Docker (sem volume — os dados duram só enquanto
# o script estiver rodando) e conecta um ao outro. Ao encerrar (Ctrl+C),
# os containers são removidos automaticamente. Uso: ./start.sh

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DIR="$ROOT_DIR/api"
MONGO_DOCKERFILE="$ROOT_DIR/docker/mongo/Dockerfile"
API_DOCKERFILE="$ROOT_DIR/docker/api/Dockerfile"

NETWORK="sipub-net"

MONGO_IMAGE="sipub-mongo:local"
MONGO_CONTAINER="sipub-mongo"
MONGO_PORT="27017"

API_IMAGE="sipub-api:local"
API_CONTAINER="sipub-api"
API_PORT="8080"

export MONGO_DB="${MONGO_DB:-sipub}"
export MONGO_COLLECTION="${MONGO_COLLECTION:-movies}"
MONGO_URI="${MONGO_URI:-mongodb://${MONGO_CONTAINER}:${MONGO_PORT}}"

cleanup() {
  echo
  echo "==> Encerrando e removendo os containers..."
  docker rm -f "$API_CONTAINER" "$MONGO_CONTAINER" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "==> Criando rede Docker (se ainda não existir)..."
docker network inspect "$NETWORK" >/dev/null 2>&1 || docker network create "$NETWORK" >/dev/null

echo "==> Construindo imagem do MongoDB..."
docker build -t "$MONGO_IMAGE" -f "$MONGO_DOCKERFILE" "$ROOT_DIR"

echo "==> Subindo container novo do MongoDB..."
docker rm -f "$MONGO_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$MONGO_CONTAINER" --network "$NETWORK" -p "${MONGO_PORT}:27017" "$MONGO_IMAGE" >/dev/null
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

echo "==> Construindo imagem da API..."
docker build -t "$API_IMAGE" -f "$API_DOCKERFILE" "$API_DIR"

echo "==> Subindo a API Go em http://localhost:${API_PORT} (Ctrl+C para encerrar)..."
docker rm -f "$API_CONTAINER" >/dev/null 2>&1 || true
docker run --name "$API_CONTAINER" --network "$NETWORK" \
  -p "${API_PORT}:8080" \
  -e "MONGO_URI=${MONGO_URI}" \
  -e "MONGO_DB=${MONGO_DB}" \
  -e "MONGO_COLLECTION=${MONGO_COLLECTION}" \
  "$API_IMAGE"
