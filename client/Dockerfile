# Build stage
FROM node:slim AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

# Development stage
FROM node:23.11-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install --only=production

COPY . .

COPY --from=builder /app/build ./build

EXPOSE 3000

CMD ["npm", "start"]
