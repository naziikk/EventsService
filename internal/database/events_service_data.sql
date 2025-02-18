DROP SCHEMA IF EXISTS events_service CASCADE;

CREATE SCHEMA events_service;

CREATE TYPE PaymentStatus AS ENUM ('pending', 'completed', 'failed');

CREATE TYPE InviteStatus AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE events_service.events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_name VARCHAR(255) NOT NULL,
    event_description TEXT NOT NULL,
    places_count INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    organizer_id INTEGER NOT NULL REFERENCES events_service.users(id),
    venue VARCHAR(255) NOT NULL,
    is_private BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events_service_data.event_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id INTEGER NOT NULL REFERENCES events_service.events(id),
    user_id INTEGER NOT NULL REFERENCES events_service.users(id),
    status InviteStatus DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events_service.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE events_service.tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL REFERENCES events_service.users(id),
);

CREATE TABLE events_service.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES events_service.users(id),
    ticket_id INTEGER NOT NULL REFERENCES events_service.tickets(id),
    amount DECIMAL(10, 2) NOT NULL,
    payment_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    payment_status PaymentStatus DEFAULT 'pending'
);