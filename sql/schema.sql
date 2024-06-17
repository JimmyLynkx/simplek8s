-- init.sql

-- 创建 cluster 表
CREATE TABLE clusters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    config TEXT NOT NULL
);