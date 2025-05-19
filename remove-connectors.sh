#!/bin/bash

CONNECT_URL="http://localhost:8083"

connectors=$(curl -s "${CONNECT_URL}/connectors" | jq -r '.[]')

if [ -z "$connectors" ]; then
    echo "Nenhum conector encontrado para remover."
else
    echo "Removendo conectores:"
    for connector in $connectors; do
        echo "Removendo conector: $connector"
        curl -s -X DELETE "${CONNECT_URL}/connectors/${connector}"
        echo "Conector $connector removido."
    done
fi

echo
echo "Todos os conectores foram removidos:"
curl -s "${CONNECT_URL}/connectors"
echo
