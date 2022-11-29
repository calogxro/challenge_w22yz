CREATE DATABASE IF NOT EXISTS qaservice;

CREATE TABLE IF NOT EXISTS qaservice.answer (
    key_        VARCHAR(128) NOT NULL,
    value       VARCHAR(128) NOT NULL
);