#!/bin/bash

MONGO_DATABASE=${DB_NAME}
MONGO_USER=${DB_USER}
MONGO_PASSWORD=${DB_PASSWORD}
ROOT_USERNAME = ${DB_ROOT_USERNAME}
ROOT_PASSWORD = ${DB_ROOT_PASSWORD}
echo "Starting MongoDB initialization..."
mongosh -u "$ROOT_USERNAME" -p "$ROOT_PASSWORD" admin <<EOF
try {
    use $MONGO_DATABASE;
    db.createUser({
        user: "$MONGO_USER",
        pwd: "$MONGO_PASSWORD",
        roles: ["readWrite"]
    });
    print("Successfully created user and database.");
} catch (e) {
    print("Error during MongoDB initialization: " + e);
}
EOF

echo "MongoDB initialization script finished."