yakdb [![Build Status](https://travis-ci.org/needcaffeine/yakdb.svg?branch=master)](https://travis-ci.org/needcaffeine/yakdb)
=====

**yakdb** (**y**et **a**nother **k**ey-value **d**ata**b**ase) is a highly-performant in-memory key-value store written in Go.

yakdb is very simple. You can get data, put data, delete data. Any and all other operations must be handled on the application side. This allows for the fastest path to maximum stability.

# Usage
Install go, then run:

    $ go get github.com/needcaffeine/yakdb/...

This will install the yakdb binary utility into your $GOBIN path.

By default yakdb will run on port :9532 (:ykdb).

    $ ./yakdb
    Listening on port 9532...

As said previously, yakdb currently supports only getting, putting, and deleting data. Let's see how these work.

    $ curl -i http://localhost:9532
    This is yakdb, a highly performant key-value store written in Go.

    Usage:
    ------
    List all items: GET /items
    Get an item: GET /items/{itemid}
    Put an item: PUT /items
    Delete an item: DELETE /items/{itemid}
    Delete all items: DELETE /items

    More documentation: https://github.com/needcaffeine/yakdb

Since yakdb is an in-memory store, we won't start out with any data. Let's verify that.

    $ curl http://localhost:9532/items
    {}

Checks out. Let's insert some data.

    $ curl -H 'Content-Type: application/json' -XPUT http://localhost:9532/items \
      -d '{
        "id": "My favorite TV shows",
        "value": "Firefly, BSG, Star Trek, SG1"
      }'
    {"status":"OK"}

You can also add complex JSON documents.

    $ curl -H 'Content-Type: application/json' -XPUT http://localhost:9532/items \
      -d '{
        "id": "Firefly",
        "value": "{\"Name\": \"Firefly\", \"Genre\": [\"Space Western\", \"Drama\", \"Science fiction\"]}"
      }'
    {"status":"OK"}

Let's see what we have in our items list now.

    $ curl http://localhost:9532/items
    {
      Firefly: {
        id: "Firefly",
        value: "{
          "Name": "Firefly",
          "Genre": ["Space Western", "Drama", "Science fiction"]
        }",
        age: 19
      },
      My favorite TV shows: {
        id: "My favorite TV shows",
        value: "Firefly, BSG, Star Trek, SG1",
        age: 50
      }
    }

The PUT operation is idempotent. If there is an existing item with the provided key, it will be overwritten. Else, created.

Now for the last two operations.

    $ curl -XDELETE http://localhost:9532/items/Firefly
    {"status":"OK"}

    $ curl http://localhost:9532/items
    {
      My favorite TV shows: {
        id: "My favorite TV shows",
        value: "Firefly, BSG, Star Trek, SG1",
        age: 50
      }
    }

Well that worked. And finally, a flush operation.

    $ curl -XDELETE http://localhost:9532/items
    {"status":"OK"}

    $ curl http://localhost:9532/items
    {}

Bam, clean slate.
