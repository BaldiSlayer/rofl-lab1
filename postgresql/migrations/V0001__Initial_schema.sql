CREATE SCHEMA tfllab1;

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
