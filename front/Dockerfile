FROM node:16-alpine as base

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package*.json ./
RUN npm ci

FROM base as dev
CMD ["npm", "run", "dev-docker"]

FROM base as prod
COPY . .
CMD ls -l /app && npm run build && ls -l dist