# google-groups

Recursively fetch GSuite groups by user email given a service account and impersonation email

## Installation

```
go get -u github.com/onetwopunch/google-groups
```

## Usage

```
Usage of google-groups:
  -depth int
    	Depth of recursion. i.e user belongs to group, which belongs to group, etc
  -impersonate string
    	GSuite admin email to impersonate
  -key-file string
    	Service Account Key JSON file path
  -subject string
    	GSuite user for which to fetch groups
```
