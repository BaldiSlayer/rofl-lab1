CREATE TABLE tfllab1.user_lock (
       user_id BIGINT PRIMARY KEY,
       expires_at TIMESTAMP NOT NULL
);

CREATE INDEX CONCURRENTLY user_lock_expires_at_index ON tfllab1.user_lock (expires_at);
CREATE INDEX CONCURRENTLY user_lock_instance_id_index ON tfllab1.user_lock (instance_id);
