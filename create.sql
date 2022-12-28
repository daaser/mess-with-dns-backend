CREATE TABLE IF NOT EXISTS dns_serials
(
    serial INT
);

CREATE TABLE IF NOT EXISTS subdomains
(
    name TEXT
);

CREATE TABLE IF NOT EXISTS dns_requests
(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT,
    subdomain TEXT,
    request TEXT,
    response TEXT,
    src_ip TEXT,
    src_host TEXT
);

CREATE TABLE IF NOT EXISTS dns_records
(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT,
    subdomain TEXT,
    rrtype TEXT,
    content TEXT
);
