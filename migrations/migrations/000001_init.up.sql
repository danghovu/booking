CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    available_seats INTEGER NOT NULL,
    start_at TIMESTAMP NOT NULL,
    location VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    price BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    creator_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    initial_quantity INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_bookings_event FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE booking_items (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_booking_items_booking FOREIGN KEY (booking_id) REFERENCES bookings(id)
);

CREATE TABLE event_tokens (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    token VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    holder_id INTEGER,
    locked_until TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_event_tokens_event FOREIGN KEY (event_id) REFERENCES events(id)
);

INSERT INTO users (email, password, role) VALUES ('admin@example.com', '$2a$10$W.klY/GB4T1EsLw8gLZI3u/.EgsYtZiCfOKIRQeGN9/V17frLc1vC', 'admin');