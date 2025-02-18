DROP SCHEMA IF EXISTS events_service_data CASCADE;

CREATE SCHEMA events_service_data;

CREATE TABLE events_service_data.events (
    id SERIAL PRIMARY KEY,
    event_name VARCHAR(255) NOT NULL,
    places_count INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events_service.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    budget DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE events_service_data.tickets (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES events_service.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);