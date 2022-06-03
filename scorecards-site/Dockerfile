FROM node:14

WORKDIR /usr/src/app

EXPOSE 8080
ENV VUE_APP_FRONTEND=http://0.0.0.0
ENV HOST=0.0.0.0
ENV PORT=8080

COPY . ./
RUN yarn install

RUN yarn build

CMD [ "yarn", "start" ]
