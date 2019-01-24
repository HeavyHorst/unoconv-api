# unoconv-api

unoconv-api is a simple http service that will call the unoconv executable, that in turn will call into the libreoffice api to convert documents.

For more info on unoconv and all possible conversions see the [unoconv website](http://dag.wiee.rs/home-made/unoconv/).

Note: unoconv-api will serialize calls to the unoconv executable, since it doesn't support concurency.


## Usage

Post the file you want to convert to the server and get the converted file in return.

The API for the webservice is /unoconv/{format-to-convert-to} so a docx to pdf conversion would be:

```sh
$ curl --form file=@myfile.docx http://127.0.0.1:3000/unoconv/pdf > myfile.pdf
```

There's a simple healthcheck api on `curl http://127.0.0.1:3000/unoconv/health` that tries to convert some plain text to pdf,
and returns "200 -OK" on success.

To compile just run `go build`, the project uses [go 1.11 modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more),
so it'll download dependencies automatically.


## systemd services

Two services, one for unoconv (the libreoffice listener), and the other for unoconv-api.
Both are running with the same user, but also DynamicUser=true, so no need to create the user
on the system. They share a runtime directory (`/run/unoconv`) for temporary files.

The services are bound together (`BindsTo=`), so if one is restarted the other will be too.
unoconv-api implements a systemd watchdog too, that will test that conversions are possible.

```
# unoconv-api.service
[Unit]
Description=unoconv-api listener
BindsTo=unoconv.service
After=unoconv.service

[Service]
Type=notify
DynamicUser=yes
WatchdogSec=60

User=unoconv
Group=unoconv
RuntimeDirectory=unoconv
RuntimeDirectoryPreserve=yes

Environment=TMPDIR=/run/unoconv
Environment=LISTEN_ADDR=127.0.0.1:3000
ExecStart=/usr/local/bin/unoconv-api
SyslogIdentifier=unoconv-api
Restart=always

[Install]
WantedBy=multi-user.target
```

```
# unoconv.service
[Unit]
Description=unoconv libreoffice listener

[Service]
Type=simple
DynamicUser=yes

User=unoconv
Group=unoconv

RuntimeDirectory=unoconv
RuntimeDirectoryPreserve=yes
StateDirectory=unoconv

Environment=HOME=/var/lib/private/unoconv/
ExecStart=/usr/bin/unoconv --listener --server=127.0.0.1 --port=2002
SyslogIdentifier=unoconv
Restart=always

[Install]
WantedBy=multi-user.target
```
