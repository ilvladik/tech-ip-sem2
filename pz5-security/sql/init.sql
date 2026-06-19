CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    study_group TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

INSERT INTO students (full_name, study_group, email)
VALUES
    ('Иванов Иван Иванович', 'ИТТ-01-25', 'ivanov@example.com'),
    ('Петрова Мария Сергеевна', 'ИТТ-02-25', 'petrova@example.com'),
    ('Сидоров Алексей Дмитриевич', 'ИТТ-03-25', 'sidorov@example.com')
ON CONFLICT (email) DO NOTHING;
