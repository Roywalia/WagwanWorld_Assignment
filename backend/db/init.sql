DROP TABLE IF EXISTS guests;
DROP TABLE IF EXISTS events;

-- Create events table
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMP,
    location VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add a few sample guests
INSERT INTO guests (event_id, name, email, phone, status, plus_ones, dietary_restrictions) VALUES
(1, 'Sarah Johnson', 'sarah@example.com', '+1234567890', 'attending', 1, 'Vegetarian'),
(1, 'Mike Chen', 'mike@example.com', '+1987654321', 'pending', 0, ''),
(2, 'Emma Wilson', 'emma@example.com', '+1122334455', 'declined', 0, 'Gluten-free'),
(3, 'James Lee', 'james@example.com', '+1555666777', 'attending', 2, '');