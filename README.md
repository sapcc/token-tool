## Token Tool

[![Build Status](https://travis-ci.org/sapcc/token-tool.svg?branch=master)](https://travis-ci.org/sapcc/token-tool)

Tool for Keystone authentication tokens.

### Usage

```
token-tool token -h
Get token from Keystone and print

Usage:
  token-tool token [flags]

Flags:
      --format string                text, json, curlrc (Default: text)
  -h, --help                         help for token
      --keystone-endpoint string     Keystone endpoint
      --password string              Password
      --project string               Project
      --project-domain-name string   Project Domain Name
      --user string                  Username
      --user-domain-name string      User Domain Name
```