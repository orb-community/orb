# Stage 0, based on Node.js, to build and compile Angular
FROM node:14.15.4 as node
WORKDIR /app
COPY package.json /app/
# I'm not sure if yarn is available from this node it imported
# but it could be installed with npm install -g yarn
# RUN npm install -g yarn
RUN yarn install
COPY ./ /app/
RUN yarn build:prod

# Stage 1, based on Nginx, to have only the compiled app, ready for production with Nginx
FROM nginx:1.13-alpine
COPY --from=node /app/dist/ /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
