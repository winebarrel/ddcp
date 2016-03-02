# ddcp

Parallel file copy command using [dd](https://en.wikipedia.org/wiki/Dd_%28Unix%29).

## Usage

```
Usage of ddcp:
  -d string
      dest
  -n int
      chunk size [mb] (default 100)
  -p  preserve attributes
  -s string
      source
```

```sh
ddcp -s source_file.dat -d dest_file.dat
```
