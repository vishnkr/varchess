FROM node:lts-alpine
RUN npm install -g http-server
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
CMD ["npm", "run" ,"serve","--port","8080"]
EXPOSE 8080
