#! /bin/bash

echo "Computing hashed password from cleartext one for user creation"
HASHED_PASSWORD=`mysql -s -u root -p"$MYSQL_ROOT_PASSWORD" --host=127.0.0.1 --port=3306 -e "select password('"$MYSQL_GOVWA_PASSWORD"');" | tail`
sed 's/PASSWORDPLACEHOLDER/'$HASHED_PASSWORD'/g' user_creation.sql | mysql --host=127.0.0.1 --port=3306 -u root -p"$MYSQL_ROOT_PASSWORD"
if [[ "$?" -eq "0" ]]; then
    echo "User creation: OK"
    echo "Importing database"
    mysql --host=127.0.0.1 --port=3306 -u root -p"$MYSQL_ROOT_PASSWORD" < create_users_table.sql
    if [[ "$?" -eq "0" ]]; then
        echo "Database import: OK"
        echo "Exiting (0)"
        exit 0
    else
        echo "Database import: KO"
        echo "Exiting (2)"
        exit 2
    fi
else
    echo "User creation: KO"
    echo "Exiting (1)"
    exit 1
fi