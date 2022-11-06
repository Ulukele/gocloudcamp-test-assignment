
echo "BUILDING..."
build_res=$(docker compose build)
if [ $? -eq 0 ]; then
    echo "BUILD SUCCESS"
else
    echo "BUILD FAILED:"
    echo "$build_res"
fi

echo "START..."
up_res=$(docker compose up -d)
if [ $? -eq 0 ]; then
    echo "START NORMALLY"
else
    echo "START FAILED:"
    echo "$up_res"
fi

echo "WAIT 3 seconds..."
sleep 3

create_res=$(curl -X POST "localhost:8080/config/" -d '{"service":"test", "data":{"key1":"value1", "key2":"value2"}}')
echo "creare results: $create_res"
echo ""

get_res=$(curl -X GET "localhost:8080/config/?service=test")
echo "get last version results: $get_res"
echo ""

upd_res=$(curl -X PUT "localhost:8080/config/" -d '{"service":"test", "data":{"key1":"value1"}}')
echo "upd config results: $upd_res"
echo ""

get_res=$(curl -X GET "localhost:8080/config/?service=test")
echo "get last version results: $get_res"
echo ""

delete_res=$(curl -X DELETE "localhost:8080/config/?service=test")
echo "delete config results: $delete_res"
echo ""

echo "DOWN SERVICE"
docker compose stop
docker compose down

