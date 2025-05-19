#!/bin/bash

CONNECT_URL="http://localhost:8083"
CONNECTORS_DIR="./connectors"

echo "Esperando o Kafka Connect iniciar..."
until curl -s "${CONNECT_URL}/connectors" >/dev/null; do
    echo "Kafka Connect ainda não está pronto. Tentando novamente em 5 segundos..."
    sleep 5
done

echo "Kafka Connect pronto. Registrando conectores..."

for file in ${CONNECTORS_DIR}/*.json; do
    CONNECTOR_NAME=$(basename "${file}" .json)

    if curl -s "${CONNECT_URL}/connectors/${CONNECTOR_NAME}" | grep -q '"name"'; then
        echo "Conector ${CONNECTOR_NAME} já registrado. Ignorando..."
    else
        echo "Registrando conector: ${CONNECTOR_NAME}"
        curl -i -X POST \
            -H "Accept:application/json" \
            -H "Content-Type:application/json" \
            --data "@${file}" \
            "${CONNECT_URL}/connectors/"
        echo "Conector ${CONNECTOR_NAME} registrado."
    fi
done
echo
echo "Todos os conectores foram registrados com sucesso:"
curl -s "${CONNECT_URL}/connectors"
echo
