# Stage 0, based on orb-ui-module, to build and compile Angular
FROM orbcommunity/orb-ui-modules as built-module

# ARG variables which direct the UI build
# can be overwritten with --build-arg docker flag
ARG ENV=""
ARG ENV_GTAGID=""

COPY ./ /app/

RUN GTAGID=$ENV_GTAGID npm run build:prod

# Stage 1, based on Nginx, to have only the compiled app, ready for production with Nginx
FROM nginx:1.13-alpine
COPY --from=built-module /app/dist/ /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
