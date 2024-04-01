# collectr

collectr is a server with single POST json endpoint `/collect` that appends
sent data into file.

When using dabase is overkill.

**Use-cases**

- saving emails from product launch marketing sites
- saving simple contact forms from websites


## conf

The `collectr.json` file must be present before launching the server.

The `"file"` key defines where should collected data be stored.

The `"fields"` defines what structure do we expect to store and what validation rules to apply.

See [collectr_example.json](./collectr_example.json) to see the structure.

<details>
<summary><strong>CLI options</strong></summary>

<pre>
-logfile string
  path to app log file (default "json.log")
-port int
  server port number (default 8000)
</pre>
</details>

## dev

```sh
go run cmd/collectr/collectr.go
```

## prod

```sh
go build cmd/collectr/collectr.go
./collectr
```
