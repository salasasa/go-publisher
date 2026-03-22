数据持久化：
mkdir -p ~/mysql_data

启动容器：
docker run -d \
  --name mysql_container \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=your_strong_password \
  -v ~/mysql_data:/var/lib/mysql \
  mysql:latest

添加--restart unless-stopped 参数，使容器自启动：