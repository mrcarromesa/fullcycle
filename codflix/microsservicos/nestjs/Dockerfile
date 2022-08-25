FROM node:14.15.4-slim

# usuario do container = root
USER root

# Certificates
ADD cert.crt /usr/local/share/ca-certificates/cert.crt
RUN chmod 644 /usr/local/share/ca-certificates/cert.crt
RUN cat /usr/local/share/ca-certificates/cert.crt >> /etc/ssl/certs/ca-certificates.crt

WORKDIR /home/node/app

# Irá rodar o comando npm install
# Irá ler os arquivos da pasta null que não tem nada só para o terminal ficar travado ali
CMD ["sh", "-c", "npm config set strict-ssl false && npm config set registry http://registry.npmjs.org/ && npm install && tail -f /dev/null"]