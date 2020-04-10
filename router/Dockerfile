FROM alexellis2/nginx-arm:latest
WORKDIR /etc/nginx
COPY ./router/nginx.conf ./nginx.conf
CMD ["nginx", "-g", "daemon off;"]
