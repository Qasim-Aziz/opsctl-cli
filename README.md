# opsctl 

### cpu stress testing

i have created this tool for stress testing (most of the cases when you want to test autoscalling we need to generate fake load). i hope 
will extend its functionality to integrate more tools for testing puposes. 

### api-stress example

```
./bin/opsctl api-stress -c 5 -u http://lasaccounting.178.128.141.35.nip.io -e /api/accountBalance/DeductBalance -m POST -H "Content-Type:application/json,Authorization:Bearer YourToken" -b '{ "isCredit": true, "totalPrice": 0, "availableBalancec": 0, "pinCode": "string", "currency": "string", "agencyId": "string", "bookingTransactionId": "string", "transactionId": 0, "remarks": "string", "supplierCode": "string", "supplierAmount": 0, "tranDate": "2023-12-31T16:53:31.062Z", "due_date": "2023-12-31T16:53:31.062Z", "description": "string", "reference": "string", "markup": 0, "agencyMarkup": 0, "max": 0, "due": 0, "received": 0, "passenger": "string", "totalPax": 0, "invoiceType": "string", "fee": 0, "subTotal": 0, "discount": 0, "createdBy": "string", "invStatus": "string", "transactionStatus": "string" }' -O output.json -d 2s -
```

### dns-check example

```
./bin/opsctl dns-check -H dev.tripon.io --debug

```