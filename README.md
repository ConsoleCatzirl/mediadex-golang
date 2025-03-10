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
A complete list of configuration options:
```yaml
mediadex:
  paths:
    music:
      - /path/to/music
    movies:
      - /path/to/some/movies
      - /path/to/more/movies
    series:
      - /path/to/tv/shows
  backend:
    arangodb:
      auth:
        host: arangodb.example.net
        port: 9001
        user: arangoUser
        pass: arangoPassword
        insecure: true
      collections:
        database_name: mediadex-test
        collection-prefix: test-
        movie_collection: movies
        music_collection: music
        series_collection: series
    opensearch:
      auth:
        hosts:
          - node.example.net:9201
          - node.example.net:9202
        user: opensearchUser
        pass: opensearchPassword
        insecure: true
      index:
        index_prefix: mediadex-test-
        movie_index: movies
        music_index: music
        series_index: series
        replicas: 1
        shards: 2
```

## TODO

* Real documentation
* Skip items that have already been indexed
* Vaccuum (remove) non-existent files from backends
* Get metadata from various websites
