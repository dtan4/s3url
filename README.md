# s3url(1)

Retrive S3 object pre-signed URL

## How to install

TBD

## How to use

You need to set AWS credentials beforehand.

```bash
export AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_REGION=xx-yyyy-0
```

Just type the command below and get Pre-signed URL on the screen.

```bash
$ s3url -b BUCKET -k KEY [-d DURATION]
```

### Options

|Option|Description|Required|Default|
|---------|-----------|-------|-------|
|`-b`, `-bucket`|Bucket name|Required||
|`-k`, `-key`|Object key|Required||
|`-d`, `-duration`|Valid duration in minutes||5|

## Development

Retrive this repository and build using `make`.

```bash
$ go get -d github.com/dtan4/s3url
$ cd $GOPATH/src/github.com/dtan4/s3url
$ make
```

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
