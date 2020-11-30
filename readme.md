# myRestDemo
Learning go restful api and simple demo

## server/httpserver.go
    the main server code
    can get logs from localhost:8000/log.txt
    can get and response request from client/httpclient
    
## server/dohandle/handler.go
    the handler of server code
    
## client/httpclient.go
    the client, can get/update/post message to server
    Usage of ./httpclient:
      -author string
            input book author (default "Duanhong Jian")
      -id string
            input book id (default "1")
      -isbn string
            input book Isbn (default "110110110")
      -method string
            ***methods***:
             -method GETALL
            GET, for example: -method GET -id 2
            POST, for example: -method POST -name 'who are you' -isbn 19961124 -author duanhongjian 
            DELETE, for example: -method DELETE -id 1
            UPDATE, for example: -method UPDATE -id 2 -name 'who are you' -isbn 19961124 -author duanhongjian 
             (default "GETALL")
      -name string
            input book name (default "Book One")
            
## TODO
    close gracefully:use channel and http package
    error process: ......
    furthermore: High Concurrency &  high availability & lock
