See `DESIGN.md` for architecture and design details.

# Developing

## Dependencies

* Golang runtime (built on 1.11, may work with previous versions).
* etcd (provided Docker Compose setup suggested).

## Building

`./build.sh` to fetch Go dependencies and build k2.

## Running

### Using docker-compose etcd
`docker-compose up`

### k2

`./k2`
