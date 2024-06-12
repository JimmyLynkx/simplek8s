-- init.sql

-- 创建 cluster 表
CREATE TABLE clusters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    server VARCHAR(255) NOT NULL,
    config TEXT NOT NULL
);