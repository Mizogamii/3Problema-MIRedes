# Problema 3 do MI de Concorrência e Conectividade.

Projeto desenvolvido para a disciplina **MI de Concorrência e Conectividade (TEC502)**.

## Pré-requisitos
- Go >= 1.25  
- Docker  
- Ganache UI
- Make

> ⚠️ Abra o terminal na **pasta raiz do projeto** antes de rodar qualquer comando.

### 1. Inicie a blockchain
> use o quickstart do ganache

## 2. Configurando os contratos

1. Copie a chave de uma das contas padrão do ganache
2. Coloque a chave no **ADMIN_KEY** no deploy.go
3. Execute:

```bash
make deploy-contract
```
para fazer deploy nos contratos

3. Copie os endereços dos contratos informados pelo deploy

## 3. Configure os endereços e chaves no Makefile

1. Modifique os endereços de contrato para os mostrados com o deploy
```bash
CONTRACT_ADDR_REGISTRO
CONTRACT_ADDR_HISTORICO
```
2. Copie as chaves de 3 contas diferentes do ganache e cole elas nos espaços
```bash
KEY_SERVER_1
KEY_SERVER_2
KEY_SERVER_3
```

## 4. Execute os servidores 

em cmds diferentes execute:
```bash
make run-pair1
make run-pair2
make run-pair3
```

## 5. Execute 1 ou mais clientes

para cada cliente, abra um cmd e use o comando:
```bash
make run-client
```
