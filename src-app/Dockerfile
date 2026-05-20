FROM node:20-alpine 

RUN apk update && apk upgrade  

COPY package*.json ./

RUN npm install --only=production

COPY server.js .

EXPOSE 80

CMD ["npm", "start"]
