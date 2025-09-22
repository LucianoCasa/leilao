# Leilão

## Como rodar

```cmd
docker compose up --build

#DOWN
docker compose down
``` 
O serviço de leilão estará disponível em `http://localhost:8080`.

## Testes

* Criar leilão:
```cmd
curl -X POST "http://localhost:8080/auction" -H "Content-Type: application/json" -d "{\"product_name\":\"PlayStation 5\",\"category\":\"eletronics\",\"description\":\"Console PS5 seminovo com 2 controles\",\"condition\":0}"

```

* Lances:

```cmd
# Primeiro lance:
curl -X POST http://localhost:8080/bid -H "Content-Type: application/json" -d '{\"auction_id\": \"c6ee0752-64b2-4e20-b8ef-696c205f6e44\", \"user_id\": \"11111111-2222-3333-4444-555555555555\", \"amount\": 1000 }'

# Segundo lance (maior que o primeiro):
curl -X POST http://localhost:8080/bid -H "Content-Type: application/json" -d '{\"auction_id\": \"c6ee0752-64b2-4e20-b8ef-696c205f6e44\", \"user_id\": \"22222222-4444-5555-6666-777777777777\", \"amount\": 1500 }'
```

* Status do leilão:

```cmd
curl -X GET "http://localhost:8080/auction/c6ee0752-64b2-4e20-b8ef-696c205f6e44"
```

* Consultar o vencedor:
```cmd
curl -X GET "http://localhost:8080/auction/winner/0bc6c570-83ee-4f65-828e-12dcc83e9218"
```

* Testar o fechamento automático do leilão (após o tempo configurado em `AUCTION_DURATION` no arquivo `.env`):

1) Criar um leilão.
2) Aguardar o tempo configurado.
3) Consultar o status do leilão (deve estar como "status" : 1).

## Teste automatizado

```cmd
docker compose exec app go test ./internal/infra/database/auction -run TestAutoCloseAuction -v
```
