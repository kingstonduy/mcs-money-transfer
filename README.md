curl --location 'localhost:7201/api/v1/moneytransfer' \
--header 'Content-Type: application/json' \
--data '{
    "fromAccountId": "OCB12345",
    "toAccountId": "TMCP23456", 
    "amount": 1
}'