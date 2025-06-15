CREATE TABLE IF NOT EXISTS projects (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goods (
                       id SERIAL PRIMARY KEY,
                       project_id INTEGER REFERENCES projects (id),
                       name VARCHAR(255) NOT NULL,
                       description VARCHAR(255) DEFAULT '',
                       priority INTEGER DEFAULT 1,
                       removed BOOL NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX good_id_idx ON goods (id);
CREATE INDEX good_project_id_idx ON goods (project_id);
CREATE INDEX good_name_idx ON goods (name);

CREATE INDEX project_id_index ON projects (id);

INSERT INTO projects (id, name) VALUES (1, 'Первая запись');

CREATE OR REPLACE FUNCTION set_priority_on_insert() RETURNS trigger AS $$
DECLARE
    max_priority INTEGER;
BEGIN
    SELECT COALESCE(MAX(priority), 0) INTO max_priority
    FROM goods
    WHERE project_id = NEW.project_id;

    IF NEW.priority IS NULL OR NEW.priority = 0 OR NEW.priority <= max_priority THEN
        NEW.priority := max_priority + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_priority_on_insert
    BEFORE INSERT ON goods
    FOR EACH ROW
EXECUTE PROCEDURE set_priority_on_insert();