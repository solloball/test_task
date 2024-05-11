# test_task

Run from command line:
- go build -o app
- ./app

Run from docker:
- docker build --tag app .   
- docker run --publish 6060:8080 --name test --rm app 
