# Mediadex

Index metadata about media libraries into ArangoDB and/or OpenSearch

## Usage

### Quick Start

```
mediadex --config /path/to/config.yaml
```

### Config File

A minimal config file contains a search path and authentication for a backend.

#### Minimal ArangoDB Example
An example that indexes a single directory of movies into an ArangoDB database
listening on localhost port 8259:
```yaml
mediadex:
  paths:
    movies:
     - /path/to/movies
  backend:
    arangodb:
      auth:
        user: someUser
        pass: secretPassword
```

#### Minimal OpenSearch Example
An example indexing a single directory of music into an OpenSearch cluster
listening on localhost port 9200:
```yaml
mediadex:
  paths:
    music:
     - /path/to/music
  backend:
    opensearch:
      auth:
        user: someUser
        pass: secretPassword
```

#### Full Example

Both of the minimal examples assume that the backend is running on `localhost`,
but alternate hosts and ports can be specified.

There are also settings for index and collection naming, and for tuning channels
and go routines.

A complete list of configuration options:
```yaml
mediadex:
  paths:
    movies:
      - /path/one
      - /path/two
    music:
      - /path/three
    series:
      - /path/four
      - /path/five
      - /path/six
  backends:
    arangodb:
      auth:
        insecure: false
        host: "arangodb.example.com"
        port: 8529
        user: service-user
        pass: service-pass
      backend:
        database_name: "cool-project"
        collection_prefix: "my-project-"
        collection_suffix: "-scope0"
        movies_collection: "Movies"
        music_collection: "Music"
        series_collection: "Series"
    opensearch:
      auth:
        insecure: false
        hosts:
          - "https://node1.example.com:9200"
          - "https://node2.example.com:9200"
          - "https://node2.example.com:9201"
        user: service-user
        pass: service-pass
      backend:
        index_prefix: "cool-project-"
        index_suffix: "-scope0"
        movie_index: "my-Movies"
        music_index: "my-Music"
        series_index: "my-Series"
        replicas: 2
        shards: 3
  threads:
    workers_per_path_type: 16
    channel_size_files: 256
    channel_size_backend: 384
```

## TODO

* Real documentation
* Logging levels / verbosity
* Vaccuum (remove) non-existent files from backends
* Get metadata from various websites
