CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100)  NOT NULL,
    email           VARCHAR(100)  UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    role            VARCHAR(20) NOT NULL,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TABLE rooms (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100)  NOT NULL,
    capacity        INT NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TABLE bookings (
    id              SERIAL PRIMARY KEY,
    room_id         INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(150) NOT NULL,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    status          VARCHAR(20) NOT NULL,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

INSERT INTO users (name, email, password_hash, role) VALUES
('Administrator', 'admin@test.com', '$2y$10$.vx8ZWk4U0uf1WEZdxlPGONM6yOROwzTlHk/F6liK7nk9Jm5JVgHu', 'ADMIN'),
('User Satu', 'user1@test.com', '$2y$10$.vx8ZWk4U0uf1WEZdxlPGONM6yOROwzTlHk/F6liK7nk9Jm5JVgHu', 'USER'),
('User Dua', 'user2@test.com', '$2y$10$.vx8ZWk4U0uf1WEZdxlPGONM6yOROwzTlHk/F6liK7nk9Jm5JVgHu', 'USER');
