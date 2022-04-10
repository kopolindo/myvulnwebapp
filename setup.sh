#! /bin/bash

sed 's/PASSWORDPLACHOLDER/'"$MYSQL_GOVWA_PASSWORD"'/g' user_creation.sql | mysql -u root -p"$MYSQL_ROOT_PASSWORD"
mysql -u root -p"$MYSQL_ROOT_PASSWORD" < create_users_table.sql
