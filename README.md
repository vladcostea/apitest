#### apitest

Concurrently test multiple API endpoints

```
$> touch suite.yml

$> cat <<EOT >> suite.yml
config:
  # base URI that will be used to build each URL
  base_uri: 
  # ":" separated HTTP Basic Auth
  # ex: user:pass
  auth:
  # List of HTTP headers that will be used for each request
  headers:
    - { name: Content-Type, value: application/json }
tests:
  - desc: "Super descriptive text"
    path: /path/to/test?with=query
    json: '{"data": {}}'
    # Currently only matches on status code
    # This is the status code the test expects to receive
    status: 200
    # Per test specific HTTP headers
    # These will overwrite the suite level ones
    headers:
      - { name: Accept, value: text/html }
EOT

$> go run main.go
```

#### TODO

* Custom HTTP methods on a per suite or per test level. Only `POST` supported for now
* Read request body from files iso plaintext
* Custom output formatters
* Add verbose mode (print out the complete request)
* Allow for multiple test suites (right now `suite.yml` is hardcoded)
* Provide binary?
