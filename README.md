# s3url(1)

[![Build Status](https://travis-ci.org/dtan4/s3url.svg?branch=master)](https://travis-ci.org/dtan4/s3url)
[![codecov](https://codecov.io/gh/dtan4/s3url/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/s3url)
[![GitHub release](https://img.shields.io/github/release/dtan4/s3url.svg)](https://github.com/dtan4/s3url/releases)

Generate [S3 object pre-signed URL](http://docs.aws.amazon.com/AmazonS3/latest/dev/ShareObjectPreSignedURL.html) in one command

```bash
$ s3url s3://my-bucket/foo.key
https://my-bucket.s3-ap-northeast-1.amazonaws.com/foo.key?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA***************************%2Fap-northeast-1%2Fs3%2Faws4_request&X-Amz-Date=20160923T010227Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&X-Amz-Signature=****************************************************************
```

## Contents

* [Installation](#installation)
  + [Using Homebrew (OS X only)](#using-homebrew-os-x-only)
  + [Precompiled binary](#precompiled-binary)
  + [From source](#from-source)
* [Usage](#usage)
  + [Upload file together](#upload-file-together)
  + [Options](#options)
* [Development](#development)
* [License](#license)

## Installation

### Using Homebrew (OS X only)

Formula is available at [dtan4/homebrew-dtan4](https://github.com/dtan4/homebrew-dtan4).

```bash
$ brew tap dtan4/dtan4
$ brew install s3url
```

### Precompiled binary

Precompiled binaries for Windows, OS X, Linux are available at [Releases](https://github.com/dtan4/s3url/releases).

### From source

```bash
$ go get -d github.com/dtan4/s3url
$ cd $GOPATH/src/github.com/dtan4/s3url
$ make deps
$ make install
```

## Usage

You need to set AWS credentials beforehand, or you can also use [named profile](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-multiple-profiles) written in `~/.aws/credentials` and `~/.aws/config`.

```bash
export AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_REGION=xx-yyyy-0
```

Just type the command below and get Pre-signed URL on the screen.

```bash
# https:// URL (both virtual-hosted-style and path-style)
$ s3url https://BUCKET.s3-region.amazonaws.com/KEY [-d DURATION] [--profile PROFILE] [--upload UPLOAD]
$ s3url https://s3-region.amazonaws.com/BUCKET/KEY [-d DURATION] [--profile PROFILE] [--upload UPLOAD]

# s3:// URL
$ s3url s3://BUCKET/KEY [-d DURATION] [--profile PROFILE] [--upload UPLOAD]

# Using options
$ s3url -b BUCKET -k KEY [-d DURATION] [--profile PROFILE] [--upload UPLOAD]
```

### Upload file together

If target object does not exist in the bucket yet, you can upload file with `--upload` flag before getting Pre-signed URL. Following example shows that uploading `foo.key` to `s3://my-bucket/foo.key` and getting Pre-signed URL of `s3://my-bucket/foo.key` will be executed in series.

```bash
$ s3url s3://my-bucket/foo.key --upload foo.key
uploaded: /path/to/foo.key
https://my-bucket.s3-ap-northeast-1.amazonaws.com/foo.key?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA***************************%2Fap-northeast-1%2Fs3%2Faws4_request&X-Amz-Date=20160923T010227Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&X-Amz-Signature=****************************************************************
```

### Options

|Option|Description|Required|Default|
|---------|-----------|-------|-------|
|`-b`, `-bucket=BUCKET`|Bucket name|Required (if no URL is specified)||
|`-k`, `-key=KEY`|Object key|Required (if no URL is specified)||
|`-d`, `-duration=DURATION`|Valid duration in minutes||5|
|`--profile=PROFILE`|AWS profile name|||
|`--upload=UPLOAD`|File to upload|||
|`-h`, `-help`|Print command line usage|||
|`-v`, `-version`|Print version|||

## Development

Retrieve this repository and build using `make`.

```bash
$ go get -d github.com/dtan4/s3url
$ cd $GOPATH/src/github.com/dtan4/s3url
$ make deps
$ make
```

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
