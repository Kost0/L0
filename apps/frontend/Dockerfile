FROM node:18-alpine as builder

WORKDIR /app

COPY . .

RUN mkdir -p /app/dist && cp index.html script.js /app/dist

FROM node:18-alpine

WORKDIR /app

COPY --from=builder /app/dist/ .

RUN npm install -g http-server

EXPOSE 3000

CMD ["http-server", "-p", "3000"]