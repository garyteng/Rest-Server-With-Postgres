#!/bin/bash

# Install Postgres DataBase
sudo apt-get install postgresql

# Create dataBase & Insert Fake Data
sudo -h localhost -u postgres psql -a -f createDatabase.sql
