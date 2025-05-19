# Etapa 1: Build da aplicação
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copia o código fonte da aplicação para dentro do contêiner
COPY . .

# Baixa as dependências e compila a aplicação
RUN go mod tidy
RUN go build -o consumer .

# Etapa 2: Criar o contêiner final
FROM debian

WORKDIR /root/

# Copia o binário compilado do contêiner de build
COPY --from=build /app/consumer .

# Comando para rodar a aplicação
CMD ["./consumer"]