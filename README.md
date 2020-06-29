# Play API

Microservice parsing play services to act as API

## Installation

### Clone repo

``` bash
git clone https://github.com/egeback/playapi.git
```

## Deployment options

The microservice can be deployed as standlone application or in a docker container

### Standalone golang application
  
#### Run build script in from root director

``` bash
./cmd/build.sh
```

#### Run application

``` bash
./playapi
```

### Docker container

#### Configure Docker Container

Update Dockefile (update ports)

#### 1. Using docker-compose ([link](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwi06f-GpafqAhXLo4sKHVWeA3UQFjAAegQIBBAC&url=https%3A%2F%2Fdocs.docker.com%2Fcompose%2F&usg=AOvVaw02oes91geDSZ-H__u_XMxc))

``` bash
docker-compose up -d --no-deps --build
```

#### 2. Using docker build

``` bash
docker build -t egeback_playapi .
```

Both options will run swag, build golang code and deploy container

## Using API

Swagger documenation available at [http://localhost:8080/api/swagger/index.html](http://localhost:8080/api/swagger/index.html)

## TODO

* [x] Paging support
* [ ] Newly added item
* [ ] Additional services
* [ ] Test cases
* [x] Fix swag from docker
* [x] Update README.md with documentation
* [x] Update code documentation
