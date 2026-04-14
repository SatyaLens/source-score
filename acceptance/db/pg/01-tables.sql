CREATE TABLE
    sources (
        name TEXT,
        score FLOAT DEFAULT 0 CHECK (
            score >= 0
            AND score <= 1
        ),
        uri_digest TEXT PRIMARY KEY,
        summary TEXT,
        tags TEXT,
        uri TEXT
    );

CREATE TABLE
    claims (
        source_uri_digest TEXT NOT NULL,
        summary TEXT,
        title TEXT,
        uri TEXT,
        uri_digest TEXT PRIMARY KEY,
        checked BOOLEAN DEFAULT FALSE,
        validity BOOLEAN DEFAULT FALSE,
        CONSTRAINT fk_source FOREIGN KEY (source_uri_digest) REFERENCES sources (uri_digest) ON DELETE CASCADE
    );

CREATE TABLE
    proofs (
        claim_uri_digest TEXT NOT NULL,
        supports_claim BOOLEAN,
        reviewed_by TEXT,
        uri TEXT,
        uri_digest TEXT PRIMARY KEY,
        CONSTRAINT fk_claim FOREIGN KEY (claim_uri_digest) REFERENCES claims (uri_digest) ON DELETE CASCADE
    );