FROM node:14.17.0-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./

COPY . .

FROM nginx:1.21.3-alpine


COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
