-- Create the database
CREATE DATABASE storedb;

--Create the user and grant priviledges
CREATE USER leroysb WITH PASSWORD 'leroysb';
GRANT ALL PRIVILEGES ON DATABASE storedb TO leroysb;
