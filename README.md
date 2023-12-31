# opsctl 

### cpu stress testing

i have created this tool for stress testing (most of the cases when you want to test autoscalling we need to generate fake load). i hope 
will extend its functionality to integrate more tools for testing puposes. 

### api-stress example

```
./bin/opsctl api-stress -c 5 -u http://example -e /api/test-endpoint -m POST -H "Content-Type:application/json,Authorization:Bearer YourToken" -b '{"Name": "mytestname" }' -O output.json -d 2s -
```

### dns-check example

```
./bin/opsctl dns-check -H dev.tripon.io --debug

```