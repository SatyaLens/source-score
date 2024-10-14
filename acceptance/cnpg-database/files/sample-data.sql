INSERT INTO
    sources (name, score, summary, tags, uri)
VALUES
    (
        'Source One',
        85,
        'A reliable source for fact-checking news.',
        'news,fact-checking,reliable',
        'https://source1.com'
    ),
    (
        'Source Two',
        90,
        'A well-known scientific publication.',
        'science,research,publication',
        'https://source2.com'
    ),
    (
        'Source Three',
        70,
        'An investigative journalism platform.',
        'journalism,investigation',
        'https://source3.com'
    );

INSERT INTO
    claims (source_uri, summary, title, uri, validity)
VALUES
    (
        'https://source1.com',
        'Claim that climate change is accelerating.',
        'Climate Change Acceleration',
        'https://claim1.com',
        TRUE
    ),
    (
        'https://source2.com',
        'Claim that new particles were discovered in particle physics.',
        'New Particle Discovery',
        'https://claim2.com',
        FALSE
    ),
    (
        'https://source3.com',
        'Claim about political corruption in a developing nation.',
        'Political Corruption Investigation',
        'https://claim3.com',
        TRUE
    );

INSERT INTO
    proofs (claim_uri, reviewed_by, uri)
VALUES
    (
        'https://claim1.com',
        'reviewer1@example.com',
        'https://proof1.com'
    ),
    (
        'https://claim2.com',
        'reviewer2@example.com,reviewer3@example.com',
        'https://proof2.com'
    ),
    (
        'https://claim3.com',
        'reviewer4@example.com',
        'https://proof3.com'
    );