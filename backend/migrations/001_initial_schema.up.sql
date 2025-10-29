-- Create guests table
CREATE TABLE guests (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'attending', 'declined')),
    notes TEXT,
    plus_ones INTEGER DEFAULT 0,
    dietary_restrictions TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(event_id, email)
);

-- Insert sample events
INSERT INTO events (title, description, event_date, location) VALUES
('Summer Gala 2025', 'Annual company celebration with dinner and awards', '2025-08-20 18:00:00', 'Downtown Grand Hotel'),
('Tech Conference', 'AI & Cloud Innovation Summit', '2025-09-15 09:00:00', 'Convention Center'),
('Team Building Retreat', 'Outdoor activities and strategy sessions', '2025-10-05 10:00:00', 'Mountain Lodge Resort');