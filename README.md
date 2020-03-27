# logger
# usage
---
### get
`go get github.com/jun_best/logger`  
### config file  
edit `logger.json`   
# features   
---
### log level:
- DEBUG
- INFO
- ERR   
### log format:  
`[2020-03-27 22:21:13 CST] [DEBUG] [file:100] log`  
### rotation type:
- time
- size  
default rotation type is time  
rotation file name like `logger.log_2020-03-27_22:18:19`
  
### rotation time type:
- MIN: do rotation every minute
- HOUR: do rotation every hour
- DAY: do rotation every day  
default rotation time type is DAY   

### rotation size type:
- 1024: bytes
- 10K: reach 10K bytes
- 10M: reach 10M bytes

# try
---
you can use unit test to run it:  
`go test`  
it will create `log/logger.log` and wait.  
after one minutes, it will do rotation and stop.


