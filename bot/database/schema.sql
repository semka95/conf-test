CREATE TABLE IF NOT EXISTS report(
    id integer PRIMARY KEY,
    url text CHECK (LENGTH(url) <= 100) UNIQUE NOT NULL,
    title text CHECK (LENGTH(title) <= 50) NOT NULL,
    starting_at datetime NOT NULL,
    duration_minutes integer NOT NULL CHECK (duration_minutes <= 720),
    reporters text CHECK (LENGTH(reporters) <= 50) NOT NULL,
    status text CHECK (status IN ('active', 'inactive')) NOT NULL
);

CREATE TABLE IF NOT EXISTS rating(
    id integer PRIMARY KEY,
    report_id integer NOT NULL,
    user_id integer NOT NULL,
    rating_type text CHECK (rating_type IN ('score', 'not_present', 'no_comments')) NOT NULL,
    content_score integer CHECK (content_score BETWEEN 1 AND 5),
    presentation_score integer CHECK (presentation_score BETWEEN 1 AND 5),
    notes text CHECK (LENGTH(notes) <= 200),
    FOREIGN KEY (report_id) REFERENCES report(id),
    FOREIGN KEY (user_id) REFERENCES user (telegram_id)
);

CREATE INDEX IF NOT EXISTS idx_rating_report_id ON rating(report_id);

CREATE INDEX IF NOT EXISTS idx_rating_user_id ON rating(user_id);

CREATE TABLE IF NOT EXISTS user (
    telegram_id integer PRIMARY KEY NOT NULL,
    id_data text CHECK (LENGTH(id_data) <= 50),
    role TEXT CHECK (ROLE IN ('admin', 'user')) NOT NULL
);

CREATE TABLE IF NOT EXISTS favorite_reports(
    user_id integer NOT NULL,
    report_id integer NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (telegram_id),
    FOREIGN KEY (report_id) REFERENCES report(id)
);

CREATE INDEX IF NOT EXISTS idx_favorite_reports_report_id ON favorite_reports(report_id);

CREATE INDEX IF NOT EXISTS idx_favorite_reports_user_id ON favorite_reports(user_id);

