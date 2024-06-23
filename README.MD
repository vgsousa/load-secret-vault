# Aplicação Go para substituição de secrets Vault Hashcorp
Esta é uma aplicação em Go que lê um arquivo chamado 'exemplo' e substitui certos valores usando o HashiCorp Vault. A aplicação cria um novo arquivo, mantendo os valores originais e trocando apenas os valores encontrados no Vault.

## Funcionalidades
- Recebe comando contendo o caminho do secret, arquivo de exemplo e Nome do arquivo de saída.
- Valida se todos os campos obrigatórios foram preenchidos.
- Gera arquivo de environment.

## Pré-requisitos
Certifique-se de ter os seguintes recursos instalados em sua máquina local:

- Go (você pode baixá-lo em golang.org)
- Docker (você pode baixá-lo em docker.com)
- Vault CLI (você pode baixá-lo em developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install)

## Instalação e Uso
1.Clone este repositório em sua máquina local:
```bash
git clone https://github.com/vgsousa/load-secret-vault.git
```

2.Navegue até o diretório do projeto:
```bash
cd load-secret-vault
```

3.Execute o Docker Compose para criação do Hashcorp Vault:
```bash
docker compose up
```

4.Configure os valores a seguir, podendo ser via CLI ou pela interface:
```bash
vault kv put secret/Go \
  DB_HASHICORP_HOST="localhost" \
  DB_HASHICORP_PASSWORD="s3cr3tp@ss" \
  DB_HASHICORP_PORT="5432" \
  DB_HASHICORP_SCHEMA="public" \
  DB_HASHICORP_USERNAME="admin"
```
5.O Vault estará disponível em http://localhost:8200.

6.Configure as Envs para executar o app em go:
```bash
export Address="http://localhost:8200"
export Token="root"
```

# Executar App
Esta aplicação executa com os dados:

- secretPath: Pacote dos Segredos.
- inputFile: Arquivo de Referência.
- outputFile: Arquivo que será criado.

Exemplo de execução:

```bash
go run main.go go .env.example .environment
```

Licença
Este projeto está licenciado sob a MIT License.