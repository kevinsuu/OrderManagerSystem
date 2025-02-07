FROM nginx:alpine

# 複製 Nginx 配置文件
COPY config/nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 8080

# 啟動 Nginx
CMD ["nginx", "-g", "daemon off;"]
