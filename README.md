# Gestion de facture website

## Binary param

### Flags

- u : db conn string, default=test
- password : password for db conn string, default=test
- db : host URI, required
- sysdba : bool for use sysdba, default=false
- h : website host url, default=localhost
- port : website port, default=8000

#### Example
"""

go run main.go -u "db2_fall2024" -password "Baddemon665!" -db "aepro.dev:1521"

"""

## Driver Oracle SQL
Godror [repo](https://github.com/godror/godror)
[Documentation](https://pkg.go.dev/github.com/godror/godror)
