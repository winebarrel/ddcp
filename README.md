# ddcp

Parallel file copy command using [dd](https://en.wikipedia.org/wiki/Dd_%28Unix%29).

## Usage

```
Usage of ddcp:
  -d string
      dest
  -n int
      chunk size (default 104857600)
  -p  preserve
  -s string
      source
```

```sh
ddcp -s source_file.dat -d dest_file.dat
```
