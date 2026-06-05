FROM node:20-alpine

RUN apk update && apk upgrade --no-cache && \
    npm install -g npm@latest

COPY package*.json ./

RUN npm ci --only=production

COPY server.js .

EXPOSE 80

CMD ["npm", "start"]
