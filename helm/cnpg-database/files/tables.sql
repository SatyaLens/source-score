CREATE TABLE
    sources (
        id BIGSERIAL,
        name TEXT,
        score SMALLINT DEFAULT 0 CHECK (
            score >= 0
            AND score <= 100
        ),
        uri_digest TEXT,
        summary TEXT,
        tags TEXT,
        uri TEXT PRIMARY KEY
    );


CREATE TABLE
    claims (
        id BIGSERIAL,
        source_uri TEXT,
        summary TEXT,
        title TEXT,
        uri TEXT PRIMARY KEY,
        validity BOOLEAN DEFAULT FALSE,
        CONSTRAINT fk_source FOREIGN KEY (source_uri) REFERENCES sources (uri) ON DELETE CASCADE
    );

CREATE TABLE
    proofs (
        id BIGSERIAL,
        claim_uri TEXT,
        reviewed_by TEXT,
        uri TEXT PRIMARY KEY,
        CONSTRAINT fk_claim FOREIGN KEY (claim_uri) REFERENCES claims (uri) ON DELETE CASCADE
    );