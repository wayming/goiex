# Stage 0, "build-stage", based on Node.js to build the frontend
FROM node:alpine as build
COPY ./goiex_front /app/goiex_front
WORKDIR /app/goiex_front
RUN npm install
RUN npm run build

# Stage 1, based on NGINX to provide a configuration to be used with react-router
FROM nginx:alpine
COPY --from=build /app/goiex_front/build /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY config/nginx.conf /etc/nginx/conf.d
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]