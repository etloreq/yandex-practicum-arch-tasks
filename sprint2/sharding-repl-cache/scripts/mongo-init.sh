#!/bin/bash

###
# Инициализируем сервер конфигурации
###

echo "Инициализируем сервер конфигурации"
docker compose exec -T configSrv mongosh --port 27017 <<EOF

rs.initiate(
  {
    _id : "config_server",
       configsvr: true,
    members: [
      { _id : 0, host : "configSrv:27017" }
    ]
  }
);
EOF


###
# Инициализируем шарды и их реплики
###

echo "Инициализируем шарды и их реплики (1)"

docker compose exec -T shard1 mongosh --port 27018 <<EOF
rs.initiate(
    {
      _id : "shard1",
      members: [
        { _id : 0, host : "shard1:27018" },
        { _id : 1, host : "shard1-slave1:27021" },
        { _id : 2, host : "shard1-slave2:27022" },
        { _id : 3, host : "shard1-slave3:27024" }
      ]
    }
)

EOF

echo "Инициализируем шарды и их реплики (2)"
docker compose exec -T shard2 mongosh --port 27019 <<EOF
rs.initiate(
    {
      _id : "shard2",
      members: [
        { _id : 0, host : "shard2:27019" },
        { _id : 1, host : "shard2-slave1:27023" },
        { _id : 2, host : "shard2-slave2:27025" },
        { _id : 3, host : "shard2-slave3:27026" }
      ]
    }
)

EOF


###
# Инициализируем роутер и наполняем данными
###

echo "Инициализируем роутер и наполняем данными"
docker compose exec -T mongos_router mongosh --port 27020 <<EOF

sh.addShard( "shard1/shard1:27018");
sh.addShard( "shard2/shard2:27019");

sh.enableSharding("somedb");
sh.shardCollection("somedb.helloDoc", { "name" : "hashed" } )

use somedb
for(var i = 0; i < 1000; i++) db.helloDoc.insert({age:i, name:"ly"+i})

EOF


###
# Проверяем сколько на первом шарде (мастер)
###
echo "Проверяем сколько на первом шарде"
docker compose exec -T shard1 mongosh --port 27018 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF

###
# Проверяем сколько на первом шарде (slave1)
###
echo "Проверяем сколько на первом шарде (slave1)"
docker compose exec -T shard1-slave1 mongosh --port 27021 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF

###
# Проверяем сколько на первом шарде (slave2)
###
echo "Проверяем сколько на первом шарде (slave2)"
docker compose exec -T shard1-slave2 mongosh --port 27022 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF

###
# Проверяем сколько на первом шарде (slave3)
###
echo "Проверяем сколько на первом шарде (slave3)"
docker compose exec -T shard1-slave3 mongosh --port 27024 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF

###
# Проверяем сколько на втором шарде
###
echo "Проверяем сколько на втором шарде"
docker compose exec -T shard2 mongosh --port 27019 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF

###
# Проверяем сколько на втором шарде (slave1)
###
echo "Проверяем сколько на втором шарде (slave1)"
docker compose exec -T shard2-slave1 mongosh --port 27023 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF

###
# Проверяем сколько на втором шарде (slave2)
###
echo "Проверяем сколько на втором шарде (slave2)"
docker compose exec -T shard2-slave2 mongosh --port 27025 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF

###
# Проверяем сколько на втором шарде (slave3)
###
echo "Проверяем сколько на втором шарде (slave3)"
docker compose exec -T shard2-slave3 mongosh --port 27026 <<EOF
use somedb
db.helloDoc.countDocuments()
EOF
