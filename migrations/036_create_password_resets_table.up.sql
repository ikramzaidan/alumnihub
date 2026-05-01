-- Create password_resets table
CREATE TABLE public.password_resets (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for password_resets
CREATE INDEX idx_password_resets_email ON public.password_resets(email);