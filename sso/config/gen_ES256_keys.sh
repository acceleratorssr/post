#!/bin/bash

openssl ecparam -name prime256v1 -genkey -noout -out private_key.pem
echo "private_key.pem get!"

openssl ec -in private_key.pem -pubout -out public_key.pem
echo "public_key.pem get!"
