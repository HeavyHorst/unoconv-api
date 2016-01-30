[![Build Status](https://drone.io/github.com/HeavyHorst/unoconv-api/status.png)](https://drone.io/github.com/HeavyHorst/unoconv-api/latest)

# unoconv-api

# Run - example

```sh
$ docker run -d -p 80:3000 --name unoconv quay.io/kaufmann_r/unoconv-api
```

# Usage

Post the file you want to convert to the server and get the converted file in return.

See all possible conversions on the [unoconv website](http://dag.wiee.rs/home-made/unoconv/).

API for the webservice is /unoconv/{format-to-convert-to} so a docx to pdf would be

```sh
$ curl --form file=@myfile.docx http://192.168.0.111/unoconv/pdf > myfile.pdf
```
