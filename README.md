# Stress Test em Go

## Execução

- docker build -t stress-test .
- docker run stress-test --url=http://google.com --requests=100 --concurrency=10
