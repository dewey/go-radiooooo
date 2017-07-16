# API

## Requests

### **GET** - /api/countries/all

#### Description
This also includes old countries like "British India"

#### CURL

```sh
curl -X GET "http://radiooooo.com/api/countries/all" \
    -H "Content-Type: application/json"
```

#### Header Parameters

- **Content-Type** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "application/json"
  ],
  "default": "application/json"
}
```

## References

