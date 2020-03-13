Beancount fava auth
> add a login page for beancount fava ui
---

### Usage
Download binary from [release](https://github.com/dongfg/fava-auth/releases) page
```
./fava-auth --username/-u <username> --password/-p <password> --server/-s <fava ip:port> [--port/-P <listen port>]
# default listen port is 9000
```

### How does it work
- check if has cookie `x-auth-fava` and valid (validate jwt and expire time)
- if not, redirect to login page
- after success login, write a cookie item
- cookie valid now, reverse proxy to fava ui
