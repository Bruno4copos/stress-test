# Load Tester CLI em Go

Este é um sistema CLI simples em Go para realizar testes de carga em um serviço web especificado pelo usuário. Ele permite configurar o número total de requisições e a quantidade de chamadas simultâneas para simular carga e coletar métricas básicas sobre a resposta do serviço.

## Funcionalidades

* **Entrada de Parâmetros via CLI:**
    * `--url`: URL do serviço a ser testado (obrigatório).
    * `--requests`: Número total de requests a serem enviados (padrão: 100).
    * `--concurrency`: Número de chamadas simultâneas (padrão: 10).
* **Execução de Teste de Carga:**
    * Realiza requisições HTTP GET para a URL especificada.
    * Controla o número de requisições simultâneas de acordo com o parâmetro `--concurrency`.
    * Garante que o número total de requisições definido em `--requests` seja executado.
* **Geração de Relatório:**
    * Apresenta um relatório ao final dos testes contendo:
        * Tempo total gasto na execução do teste.
        * Quantidade total de requests realizados.
        * Quantidade de requests com status HTTP 200.
        * Distribuição de outros códigos de status HTTP encontrados.
* **Execução via Docker:**
    * Um `Dockerfile` está incluído para facilitar a construção de uma imagem Docker e a execução da aplicação em um container.

## Como Usar

### Execução Local

1.  **Clone o repositório (se aplicável).**
2.  **Navegue até o diretório do projeto.**
3.  **Compile o código Go:**
    ```bash
    go build -o load-tester main.go
    ```
4.  **Execute o load tester com os parâmetros desejados:**
    ```bash
    ./load-tester --url=[http://seu-servico.com](http://seu-servico.com) --requests=1000 --concurrency=50
    ```
    Substitua `http://seu-servico.com`, `1000` e `50` pelos valores desejados.

### Execução com Docker

1.  **Certifique-se de ter o Docker instalado em sua máquina.**
2.  **Navegue até o diretório do projeto onde o `Dockerfile` está localizado.**
3.  **Construa a imagem Docker:**
    ```bash
    docker build -t load-tester .
    ```
4.  **Execute o container Docker passando os parâmetros como argumentos:**
    ```bash
    docker run load-tester --url=[http://google.com](http://google.com) --requests=500 --concurrency=20
    ```
    Adapte a URL, o número de requests e a concorrência conforme necessário.

## Relatório de Exemplo

Após a execução do teste, um relatório similar ao seguinte será exibido:

     ```bash
     Iniciando teste de carga para: http://google.com
     Total de requests: 1000
     Concorrência: 10
     
     --- Relatório de Teste de Carga ---
     Tempo total de execução: 1m3.2543876s
     Total de requests realizados: 1000
     Requests com status HTTP 200: 995
     Distribuição de outros códigos de status HTTP:
     503: 5
	 ```

## Próximos Passos e Melhorias

Este é um projeto básico e pode ser expandido com diversas funcionalidades, como:

* Suporte a diferentes métodos HTTP (POST, PUT, DELETE, etc.).
* Envio de payloads (corpo da requisição).
* Adição de headers personalizados.
* Medição de latência (tempo de resposta) para cada requisição e geração de estatísticas (média, mediana, percentis).
* Opção para definir um tempo limite para as requisições.
* Persistência dos resultados em um arquivo (CSV, JSON, etc.).
* Mais opções de configuração via linha de comando.
* Métricas mais detalhadas sobre a saúde do sistema testado.
* Implementação de estratégias de rampa de carga (aumentar gradualmente a concorrência).

Sinta-se à vontade para contribuir e adicionar melhorias a este projeto!