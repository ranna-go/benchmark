# benchmark

Simple tool to benchmark a ranna server instance.

## Usage

```
$ go run cmd/benchmark/main.go -h
Usage of C:\...\main.exe:
  -e string
        ranna API endpoint
  -loglevel int
        logger level (default 4)
  -n int
        number of total requests (default 5) 
  -p int
        parallel running requests (default 5)
  -pp
        pretty print results
  -snippets string
        snippets directory (default "./snippets")
```

## Benchmarks

> All benchmarks ran with a build of the state of commit [120da47](https://github.com/ranna-go/benchmark/tree/120da471ba71bdc4be1064125ff39050d98238e9).

The current demo installation *(`ranna.zekro.de`)* is currently running on a low-spec 1C 2GB VPS. On this server, two ranna instances are running behind a traefik reverse proxy. `https://public.ranna.zekro.de` is publicly available and strongly rate limited. `https://private.ranna.zekro.de` is only accessable from whitelisted addresses and has no rate limits. 

The following benchmarks were executed on this environment.

**Used Parameters**
| | |
|--|--|
| Endpoint | `https://private.ranna.zekro.de` |
| Snippets | `5` |
| # Requests | `100` |
| # Parallel Requests | `10` |

**Results**
| | |
|--|--|
| # Successful | `99` |
| % Successful | `99` |
| # Erroneous | `1` |
| % Erroneous | `1` |
| Average Request to Response Time | `7.5219s` |
| Average Execution Time | `6.4011s` |


![](https://i.imgur.com/zFjFEjS.png)

**Used Parameters**
| | |
|--|--|
| Endpoint | `https://private.ranna.zekro.de` |
| Snippets | `5` |
| # Requests | `100` |
| # Parallel Requests | `50` |

**Results**
| | |
|--|--|
| # Successful | `48` |
| % Successful | `48` |
| # Erroneous | `52` |
| % Erroneous | `52` |
| Average Request to Response Time | `14.7281s` |
| Average Execution Time | `6.7089s` |

![](https://i.imgur.com/vNojuFm.png)

---

The following benchmarks now were made on a 16C 32GB VPS.

**Used Parameters**
| | |
|--|--|
| Endpoint | `https://private.ranna.zekro.de` |
| Snippets | `5` |
| # Requests | `100` |
| # Parallel Requests | `10` |

**Results**
| | |
|--|--|
| # Successful | `100` |
| % Successful | `100` |
| # Erroneous | `0` |
| % Erroneous | `0` |
| Average Request to Response Time | `1.5096s` |
| Average Execution Time | `756ms` |

![](https://i.imgur.com/J673T0S.png)

**Used Parameters**
| | |
|--|--|
| Endpoint | `https://private.ranna.zekro.de` |
| Snippets | `5` |
| # Requests | `100` |
| # Parallel Requests | `50` |

**Results**
| | |
|--|--|
| # Successful | `100` |
| % Successful | `100` |
| # Erroneous | `0` |
| % Erroneous | `0` |
| Average Request to Response Time | `2.3727s` |
| Average Execution Time | `1.0720s` |

![](https://i.imgur.com/xC5ems4.png)

**Used Parameters**
| | |
|--|--|
| Endpoint | `https://private.ranna.zekro.de` |
| Snippets | `5` |
| # Requests | `1000` |
| # Parallel Requests | `100` |

**Results**
| | |
|--|--|
| # Successful | `915` |
| % Successful | `91.5` |
| # Erroneous | `85` |
| % Erroneous | `8.5` |
| Average Request to Response Time | `14.3220s` |
| Average Execution Time | `2.6824s` |

![](https://i.imgur.com/mhmvrRV.png)