Каталог с товарами

curl -X POST localhost:8080/v1/good -d '{"name":"name1", "unit_of_measure":"uom1", "country":"Country1"}'  
curl localhost:8080/v1/goods
curl localhost:8080/v1/good/1
curl -X PUT localhost:8080/v1/good -d '{"good":{"code":1, "name":"name2", "unit_of_measure":"uom1", "country":"Country1"}}'

// DELETE не смог вызвать, получаю ошибку {"code":12, "message":"Method Not Allowed", "details":[]}
curl -X DELETE localhost:8080/v1/good/ -d '{"code":1}'

gRPC на localhost:8081
Команды:
add name3 uom3 country3
list
get 4
update 4 name4 uom4 country4
delete 4

HTTP http://localhost:8080/v1/goods
http://localhost:8080/v1/good/1 // code == 1

swagger сгенерил /catalog/api/api.swagger.json но не добавил на сервер
