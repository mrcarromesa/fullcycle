FROM node:14.15.4-slim

# usuario do container = root
USER node

WORKDIR /home/node/app

# Irá rodar o comando npm install
# Irá ler os arquivos da pasta null que não tem nada só para o terminal ficar travado ali
CMD ["sh", "-c", "npm install && tail -f /dev/null"]