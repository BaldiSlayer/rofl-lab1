-- export $(cat .env | xargs) && docker run -it --rm --network rofl-lab1_default -v $(pwd)/postgresql/migrations/:/migrations/migrations urbica/pgmigrate -d /migrations -t latest migrate -t 1 -c "port=5432 host=postgres dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD"

CREATE TABLE tfllab1.user_state (
       user_id BIGINT PRIMARY KEY,
       state INT NOT NULL
);

CREATE TABLE tfllab1.extraction_result (
       user_id BIGINT PRIMARY KEY,
       user_request TEXT,
       formalize_result TEXT,
       parse_result JSONB,
       parse_error TEXT
);
