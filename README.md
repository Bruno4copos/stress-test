# Load Tester CLI

Uma ferramenta de linha de comando escrita em Go para realizar testes de carga (load testing) em serviços web. Permite simular um número de requisições HTTP com controle de concorrência, autenticação, envio de corpo e exportação de resultados.

## 🚀 Funcionalidades

- Testes com diferentes métodos HTTP: `GET`, `POST`, `PUT`, etc.
- Suporte a **headers personalizados**, como tokens de autenticação.
- Envio de **body** para requisições `POST`/`PUT`.
- Controle de **concorrência** e número total de requisições.
- Geração de **relatório ao final do teste**.
- Exportação dos resultados em formato **CSV**.

## 🛠️ Compilação e Execução

### Com Docker

```bash
docker build -t load-tester .
docker run --rm load-tester --url=http://localhost:8080 --requests=1000 --concurrency=20
```

### Localmente (Requer Go)

```bash
go build -o load-tester main.go
./load-tester --url=http://localhost:8080 --requests=100 --concurrency=5
```

## 📥 Parâmetros disponíveis

| Flag           | Descrição                                               | Exemplo                                      |
|----------------|---------------------------------------------------------|----------------------------------------------|
| `--url`        | URL do serviço a ser testado                            | `--url http://localhost:8080`                |
| `--requests`   | Total de requisições a serem feitas                     | `--requests 500`                             |
| `--concurrency`| Número de chamadas simultâneas                          | `--concurrency 10`                           |
| `--method`     | Método HTTP (GET, POST, PUT...)                         | `--method POST`                              |
| `--header`     | Headers personalizados (pode ser usado múltiplas vezes) | `--header "Authorization: Bearer token"`     |
| `--body`       | Corpo da requisição em texto (JSON, etc.)               | `--body '{"key":"value"}'`                   |
| `--csv`        | Exportar resultados para CSV                            | `--csv results.csv`                          |

## 📊 Relatório Gerado

Ao final do teste, será exibido:

- Tempo total da execução
- Total de requisições
- Quantidade de respostas com status 200
- Distribuição dos demais códigos HTTP

Se usar `--csv`, será criado um arquivo contendo:

```csv
Status Code,Duration (ms)
200,12.34
500,21.67
...
```

## 🧪 Exemplo de uso completo

```bash
docker run --rm load-tester \
  --url=http://localhost:8080/api \
  --requests=1000 \
  --concurrency=50 \
  --method POST \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer abc123" \
  --body '{"message":"hello"}' \
  --csv report.csv
```

---

Criado para facilitar testes de carga em serviços web. Sinta-se à vontade para contribuir!
