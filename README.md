# s3fs

[![Build Status](https://travis-ci.org/KyleBanks/s3fs.svg?branch=master)](https://travis-ci.org/KyleBanks/s3fs)

`s3fs` is an [Amazon S3](https://aws.amazon.com/s3/) client that provides a familiar interface for navigating and managing S3 buckets and objects. The goal of `s3fs` is to allow you to interact with Amazon S3 as you would your local filesystem.

![s3fs Demo](https://kylewbanks.com/images/post/s3fs-demo-github.gif)

## Installation

### Binary 

Simply go to the [releases page](https://github.com/KyleBanks/s3fs/releases) and download the appropriate binary for your environment.

### From Source

`s3fs` is written in [Go](https://golang.org/) and can be installed using `go get` if you have a working Go environment. 

```
go get -u github.com/KyleBanks/s3fs
go install github.com/KyleBanks/s3fs
```

## Usage 

Execute the `s3fs` command to launch the application, and then use the appropriate commands outlined below to navigate Amazon S3 as if it were a local filesystem.

```
s3fs
```

### cd

Changes the current working directory.

**Examples:**

```
# Move into a bucket
$ cd bucket 

# Move into a folder within a bucket
$ cd bucket/folder/subfolder

# Move back to root
$ cd /
```

### pwd

Prints the current working directory.

**Examples:**

```
$ cd bucket/folder
$ pwd
/bucket/folder
```

### ls

Lists current directory contents.

**Examples:**

```
# Print bucket list when at root.
$ ls
[B] bucket1
[B] bucket2
[B] bucket3

# Print object list when in a bucket.
$ cd bucket1
$ ls
[F] folder1/
[F] folder2/
    file1.txt

$ cd folder1
$ ls 
[F] subfolder/
    file2.txt
    file3.txt
```

### Other Commands

- `clear` clears all terminal output.
- `exit` quits `s3fs`.

# License

```
MIT License

Copyright (c) 2016 Kyle Banks

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
