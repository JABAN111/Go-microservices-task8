
CREATE TABLE db_stats (
    id SERIAL PRIMARY KEY,  
    words_total INT NOT NULL,
    words_unique INT NOT NULL,
    comics_fetched INT NOT NULL
);


CREATE TABLE service_stats (
    db_stats_id INT PRIMARY KEY,
    comics_total INT NOT NULL,
    FOREIGN KEY (db_stats_id) REFERENCES db_stats (id)
);


CREATE TABLE comics (
    id SERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    words TEXT[]
);

CREATE TABLE words_stats (
     word TEXT PRIMARY KEY,
     count INT NOT NULL DEFAULT 1
);

CREATE OR REPLACE FUNCTION add_word_stats(words_arr TEXT[])
RETURNS VOID AS
$$
BEGIN
INSERT INTO words_stats (word, count)
SELECT word, COUNT(*)
FROM unnest(words_arr) AS word
GROUP BY word
    ON CONFLICT (word) DO UPDATE
                              SET count = words_stats.count + EXCLUDED.count;
END;
$$ LANGUAGE plpgsql;