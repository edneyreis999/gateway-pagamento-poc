# Use uma imagem base que já tenha o Java instalado
FROM openjdk:17-jdk-slim

# Atualize os pacotes e instale Node.js
RUN apt-get update && \
    apt-get install -y curl gnupg && \
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs=22.5.1-1nodesource1 && \
    apt-get install -y git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Crie o usuário node e defina permissões
RUN useradd -m -d /home/node -s /bin/bash node && \
    mkdir -p /home/node/app && \
    chown -R node:node /home/node/app

# Instale o TypeScript globalmente
RUN npm install -g typescript

USER node

WORKDIR /home/node/app

CMD ["tail", "-f", "/dev/null"]
