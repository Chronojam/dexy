**Dexy**
----------
**Introduction**    

Dexy is a simple command line tool to do 3LO authentication with an external provider and give a token back.

**Usage**    

Dexy works on client_id/client_secret  pairs.
As such you'll need to generate one and authorize it with your provider. In [dex](https://github.com/coreos/dex) this might look like this:

```
staticClients:
- id: dexy
  secret: dexy-secret
  name: 'Dexy'
  redirectURIs:
  # This needs to be a local address, as dexy will start a webserver
  # for you to callback too, it will need to match the host/port in the dexy config
  - 'http://localhost.com:10000/oauth2/callback'`
```
Dexy also has its own configuration file, which it will search for in the following locations:
```
$HOME/.dexy.yaml
./.dexy.yaml
/etc/.dexy.yaml
``` 
```
auth:
  dex_host: "https://dex.mycompany.com"
  callback_host: "localhost"
  callback_port: 10111
  # This will generate a callbackurl like http://localhost:10111/oauth2/callback
  client_id: "dexy"
  client_secret: "dexy-secret"
```

It can only support a single provider at a time, if you need to change providers, you can delete the dexy token file in ~/.dexy-token.yaml

**Building**    

Pretty self explainatory but
```
make test
make build

Will drop the compiled binaries under .build/
```

**Contributing**    

All PR's are welcome, just open one aganist master
