openssl ecparam -name prime256v1 -genkey -noout -out private_key.pem
Write-Host "private_key.pem get!"

openssl ec -in private_key.pem -pubout -out public_key.pem
Write-Host "public_key.pem get!"
