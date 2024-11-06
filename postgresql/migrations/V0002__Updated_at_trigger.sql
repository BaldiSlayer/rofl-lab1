ALTER TABLE tfllab1.user_state ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT now();

CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_user_state_updated_at BEFORE UPDATE ON tfllab1.user_state
       FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
