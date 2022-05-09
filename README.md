## table-driver-tests-go

### Test

```bash
go test -v github.com/shivamMg/table-driven-tests-go/api/
```

### Mockgen

```bash
mockgen -source=api/db.go -destination=api/mock/db.go -package=mock Database
mockgen -source=api/auth.go -destination=api/mock/auth.go -package=mock Authenticator
```
