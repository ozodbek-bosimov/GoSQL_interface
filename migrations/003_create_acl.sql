CREATE TABLE IF NOT EXISTS acl (
    id BIGSERIAL PRIMARY KEY,
    src_ip TEXT NOT NULL,
    dst_ip TEXT NOT NULL,
    protocol TEXT NOT NULL CHECK (protocol IN ('tcp', 'udp', 'icmp', 'any')),
    src_port TEXT,
    dst_port TEXT,
    action TEXT NOT NULL CHECK (action IN ('permit', 'deny', 'permit+reflect')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_acl_src_ip ON acl(src_ip);
CREATE INDEX idx_acl_dst_ip ON acl(dst_ip);
