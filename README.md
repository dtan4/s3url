# s3url(1)

Retrive S3 object pre-signed URL

## Installation

TBD

## Usage

You need to set AWS credentials beforehand.

```bash
export AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_REGION=xx-yyyy-0
```

Put object into S3 bucket using Management Console or awscli.

```bash
$ aws s3 cp foo_file s3://BUCKET/KEY
```

Just type the command below and get Pre-signed URL on the screen.

```bash
# https:// URL
$ s3url https://s3-region.amazonaws.com/BUCKET/KEY [-d DURATION]

# s3:// URL
$ s3url s3://BUCKET/KEY [-d DURATION]

# Using options
$ s3url -b BUCKET -k KEY [-d DURATION]
```

### Options

|Option|Description|Required|Default|
|---------|-----------|-------|-------|
|`-b`, `-bucket`|Bucket name|Required (if no URL is specified)||
|`-k`, `-key`|Object key|Required (if no URL is specified)||
|`-d`, `-duration`|Valid duration in minutes||5|
|`-h`, `-help`|Print command line usage|||

## Development

Retrive this repository and build using `make`.

```bash
$ go get -d github.com/dtan4/s3url
$ cd $GOPATH/src/github.com/dtan4/s3url
$ make
```

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
