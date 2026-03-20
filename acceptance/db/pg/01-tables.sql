CREATE TABLE
    sources (
        name TEXT,
        score SMALLINT DEFAULT 0 CHECK (
            score >= 0
            AND score <= 100
        ),
        uri_digest TEXT PRIMARY KEY,
        summary TEXT,
        tags TEXT,
        uri TEXT
    );


CREATE TABLE
    claims (
        source_uri TEXT,
        summary TEXT,
        title TEXT,
        uri TEXT PRIMARY KEY,
        validity BOOLEAN DEFAULT FALSE,
        CONSTRAINT fk_source FOREIGN KEY (source_uri) REFERENCES sources (uri_digest) ON DELETE CASCADE
    );

CREATE TABLE
    proofs (
        claim_uri TEXT,
        reviewed_by TEXT,
        uri TEXT PRIMARY KEY,
        CONSTRAINT fk_claim FOREIGN KEY (claim_uri) REFERENCES claims (uri) ON DELETE CASCADE
    );