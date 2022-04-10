# GOVWA - GO Vulnerable Web Application

## Requirements
1. Local instance of MariaDB (scripts and code refers to `127.0.0.1:3306`)
2. Environmental variables
    * `MYSQL_ROOT_PASSWORD`
    * `MYSQL_GOVWA_PASSWORD`

## Setup
Run script `setup.sh`  
It will
1. Hash the `MYSQL_GOVWA_PASSWORD` password
2. Create `govwauser`, grantings basic privileges
3. Import a fictitious table, with most common usernames and password