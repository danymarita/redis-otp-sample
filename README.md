#Request OTP
``curl --location --request POST 'http://localhost:8001/otp/request' \
--header 'Content-Type: application/json' \
--data-raw '{
"user_id": 1
}'``
#Validate OTP
``curl --location --request POST 'http://localhost:8001/otp/validate' \
--header 'Content-Type: application/json' \
--data-raw '{
"user_id": 1,
"otp_code": 123456
}'``
#Prerequisites
1. Redis