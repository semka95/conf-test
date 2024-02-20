CREATE TABLE IF NOT EXISTS conference(
    id integer PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS report(
    id integer PRIMARY KEY,
    url text CHECK (LENGTH(url) <= 100) UNIQUE NOT NULL,
    title text CHECK (LENGTH(title) <= 50) NOT NULL,
    starting_at datetime NOT NULL,
    duration_minutes integer NOT NULL CHECK (duration_minutes <= 720),
    reporters text CHECK (LENGTH(reporters) <= 50) NOT NULL,
    conference_id integer NOT NULL,
    status text CHECK (status IN ('active', 'inactive')) NOT NULL,
    FOREIGN KEY (conference_id) REFERENCES conference(id)
);

CREATE INDEX IF NOT EXISTS idx_report_conference_id ON report(conference_id);

CREATE TABLE IF NOT EXISTS rating(
    id integer PRIMARY KEY,
    report_id integer NOT NULL,
    user_id integer NOT NULL,
    rating_type text CHECK (rating_type IN ('score', 'not_present', 'no_comments')) NOT NULL,
    content_score integer,
    presentation_score integer,
    notes text CHECK (LENGTH(notes) <= 200),
    FOREIGN KEY (report_id) REFERENCES report(id),
    FOREIGN KEY (user_id) REFERENCES USER (telegram_id)
);

CREATE INDEX IF NOT EXISTS idx_rating_report_id ON rating(report_id);

CREATE INDEX IF NOT EXISTS idx_rating_user_id ON rating(user_id);

CREATE TABLE IF NOT EXISTS USER (
    telegram_id integer PRIMARY KEY NOT NULL,
    id_data text CHECK (LENGTH(id_data) <= 50),
    role TEXT CHECK (ROLE IN ('admin', 'user')) NOT NULL
);

CREATE TABLE IF NOT EXISTS favorite_reports(
    user_id integer NOT NULL,
    report_id integer NOT NULL,
    FOREIGN KEY (user_id) REFERENCES USER (telegram_id),
    FOREIGN KEY (report_id) REFERENCES report(id)
);

CREATE INDEX IF NOT EXISTS idx_favorite_reports_report_id ON favorite_reports(report_id);

CREATE INDEX IF NOT EXISTS idx_favorite_reports_user_id ON favorite_reports(user_id);

