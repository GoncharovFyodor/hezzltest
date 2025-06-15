FROM clickhouse/clickhouse-server:23.8-alpine

# Отключаем авторизацию в конфиге
RUN echo "<clickhouse><users><default><password></password><networks><ip>::/0</ip></networks></default></users></clickhouse>" > /etc/clickhouse-server/users.d/default-user.xml

COPY ./migrations/clickhouse/*.sql /docker-entrypoint-initdb.d/
RUN chmod 644 /docker-entrypoint-initdb.d/*.sql