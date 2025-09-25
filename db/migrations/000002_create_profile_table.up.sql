CREATE TABLE profile (
    user_id INT PRIMARY KEY,
    profile_picture VARCHAR(255),
    fullname VARCHAR(255),
    phone VARCHAR(20),
    pin TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_profile_users FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);