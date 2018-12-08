# Secret Santa

## Prerequisites

Requires openssl to be installed.

## Installation

```bash
go get github.com/w-k/secret-santa
```

## Usage

Ask all secret santa participants to submit their public keys. Assuming that
a participant is a git user, he or she should have a `./ssh/id_rsa.pub` file
on their computer already. If not, they should follow [these instructions](https://confluence.atlassian.com/bitbucket/set-up-an-ssh-key-728138079.html) 
to create it.

Rename the files to each participant's name and put them into a folder. For example,
if John, Mary and Ashley are participating, we should end up with this folder:

```bash
pubkeys
├── Ashley
├── John
└── Mary
```

Where `Ashley` is the `id_rsa.pub` file received from Ashley and so on.

To demonstrate the random assignment works:

```bash
secret-santa -in=./pubkeys -demo
```

To generate the results, readable only by the owner of the corresponding private key:

```bash
secret-santa -in=./pubkeys -out=./result.zip
```

The result file will contain:

* one file per each participant containing the result (name of the person to give the present to)
* a bash script to decrypt the result
* a README file with instructions on how to use it

```bash
result
├── Ashley
├── John
├── Mary
├── README.md
└── decrypt
```
