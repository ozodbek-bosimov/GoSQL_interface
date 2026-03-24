CREATE TABLE IF NOT EXISTS interfaces (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    ip TEXT NOT NULL,
    mac TEXT NOT NULL,
    mtu INTEGER NOT NULL,
    status BOOLEAN NOT NULL DEFAULT true,
    ip_type BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_interfaces_name ON interfaces(name);
CREATE INDEX idx_interfaces_ip ON interfaces(ip);
