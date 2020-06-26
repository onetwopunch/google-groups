# google-groups

Recursively fetch GSuite groups by user email given a service account and impersonation email

## Installation

```
go get -u github.com/onetwopunch/google-groups
```

## Usage

```
Usage of google-groups:
  -impersonate string
    	GSuite admin email to impersonate
  -key-file string
    	Service Account Key JSON file path
  -subject string
    	GSuite user for which to fetch groups
  -depth int
    	(optional) Depth of recursion if desired. i.e user belongs to group, which belongs to group, etc
```

**Example usage:**

Assuming the user `tesla` belongs to the group `puppies` which subsequently belongs to the group `dogs`:

```
google-groups --impersonate admin@canty.wtf \
  --key-file ~/.config/gcloud/demo.json \
  --subject tesla@canty.wtf --depth 1
["puppies@canty.wtf","dogs@canty.wtf"]
```

This output is in JSON format so it can be easily piped to `jq` or other utilities.
