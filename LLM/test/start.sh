cp /LLM/data.yaml /LLM/test/get_similar_tests/data.yaml

gunicorn app.main:app --workers 1 --worker-class uvicorn.workers.UvicornWorker --bind=0.0.0.0:8100 --timeout 240 &

# дожидаемся момента, когда поднимется веб-север
URL="http://0.0.0.0:8100/ping"
MAX_ATTEMPTS=20
SLEEP_TIME=1

for (( i=1; i<MAX_ATTEMPTS+1; i++ ))
do
    response=$(curl -s $URL)

    if [[ "$response" == "\"OK\"" ]]; then
        break
    fi

    sleep $SLEEP_TIME
done

pytest -s /LLM/test/get_similar_tests/
