<h1 align="center">go-elasticache</h1>

<p align="center">
  <img src="https://img.shields.io/badge/Completed-100%25-green.svg?style=flat-square">
</p>

<p align="center">
  Thin abstraction over the Memcache client package <a href="https://github.com/bradfitz/gomemcache">gomemcache</a><br>
  allowing it to support <a href="https://aws.amazon.com/elasticache/">AWS ElastiCache</a> cluster nodes
</p>

## Explanation

When using the memcache client [gomemcache](https://github.com/bradfitz/gomemcache) you typically call a constructor and pass it a list of memcache nodes like so:

```go
mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
```

But when using the [AWS ElastiCache](https://aws.amazon.com/elasticache/) service you need to query a particular internal endpoint via a socket connection and manually parse out the details of the available cluster.

In Ruby it seems most Memcache clients are setup to handle this for you, but in Go that doesn't appear to be the case. Hence I've created this package as a thin abstraction layer on top of the gomemcache package.

## Example

Below is an example of how to use this package. 

To run it locally you will need the following dependencies installed and running:

- Memcache (e.g. `docker run -d -p 11211:11211 memcached`)
- [fake_elasticache](https://github.com/stevenjack/fake_elasticache) (e.g. `gem install fake_elasticache && fake_elasticache`)

```go
package main

import (
	"fmt"
	"log"

	"github.com/integralist/go-elasticache/elasticache"
)

func main() {
	mc, err := elasticache.New()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	if err := mc.Set(&elasticache.Item{Key: "foo", Value: []byte("my value")}); err != nil {
		log.Println(err.Error())
	}

	it, err := mc.Get("foo")
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Printf("%+v", it) 
  // &{Key:foo Value:[109 121 32 118 97 108 117 101] Flags:0 Expiration:0 casid:9}
}
```

> Note: when running in production make sure to set the environment variable `ELASTICACHE_ENDPOINT`

## Licence

[The MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright (c) 2016 [Mark McDonnell](http://twitter.com/integralist)
