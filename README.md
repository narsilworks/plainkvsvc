# plainkvsvc
A REST API server for PlainKV access. PlainKV stores and retrieves its values to a MySQL database.

## Basic Functions

### Get the value of a key
```bash
curl --location 'http://localhost:8080/api/sample_key' \
--header 'Authorization: Bearer <JWT here>'
```

### Set the value of a key
```bash
curl --location 'http://localhost:8080/api/sample_key' \
--header 'Content-Type: text/plain' \
--header 'Authorization: Bearer <JWT here>' \
--data 'This is the value of this key'
```
### Set the value of a key with a content type
```bash
curl --location 'http://192.168.1.37:15301/api/sample_key' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT here></JWT>' \
--data '{
    "document": "100000001",
    "title": "The Go Programming Language"
}'
```

> **Note**
> If the value being set is an image, JSON or XML string, it automatically retrieves its content type.
> When the key is retrieved, it will add the content type to the result:

### Delete the key
```bash
curl --location --request DELETE 'http://localhost:8080/api/sample_key' \
--header 'Authorization: Bearer <JWT here>'
```

## Convenience Functions for Tallying

### Get the tally of a key
For regular operation:

```bash
curl --location 'http://localhost:8080/api/tally/sample' \
--header 'Authorization: Bearer <JWT here>'
```
If you want to set the tally's initial count:

```bash
curl --location 'http://localhost:8080/api/tally/sample?offset=50' \
--header 'Authorization: Bearer <JWT here>'
```

### Increment the tally
```bash
curl --location --request PUT 'http://localhost:8080/api/incr/sample' \
--header 'Authorization: Bearer <JWT here></JWT>'
```

### Decrement the tally
```bash
curl --location --request PUT 'http://localhost:8080/api/decr/sample' \
--header 'Authorization: Bearer <JWT here>'
```

### Reset the tally
```bash
curl --location --request DELETE 'http://localhost:8080/api/tally/sample' \
--header 'Authorization: Bearer <JWT here>'
```