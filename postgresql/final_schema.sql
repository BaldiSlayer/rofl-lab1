-- export $(cat .env | xargs) && docker run -it --rm --network rofl-lab1_default -v $(pwd)/postgresql/migrations/:/migrations/migrations urbica/pgmigrate -d /migrations -t latest migrate -t 2 -c "port=5432 host=postgres dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD"

CREATE TABLE tfllab1.user_state (
       user_id BIGINT PRIMARY KEY,
       state INT NOT NULL,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE tfllab1.extraction_result (
       user_id BIGINT PRIMARY KEY,
       user_request TEXT,
       formalize_result TEXT,
       parse_result JSONB,
       parse_error TEXT
);

CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_user_state_updated_at BEFORE UPDATE ON tfllab1.user_state
       FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
