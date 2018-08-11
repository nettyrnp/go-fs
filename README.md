# Application for reading and saving records from multiple log files to DB

Features:

* Incremental reading of log files (i.e. log files can be huge)
* Graceful shutdown
* For each log file, a goroutine is launched, which monitors the addition of new lines to the file.
* The main goroutine gets the generated objects from the goroutines and saves them to MongoDB.

## How to launch the application

Start MongoDB database server.
* ...
* ...
* ...

Install the application from the Terminal:
```shell
go get github.com/nettyrnp/go-fs
```

Run the application from the Terminal or from under IntelliJ IDEA etc.:
```shell
go run main.go
```

In the Terminal you can observe the processing of the log files 'data/name1.log' and 'data/name2.log':

```shell
# Terminal:
2018/08/10 22:16:41 Reading: loaded 2 lines from file 'data/name2.log'
2018/08/10 22:16:41 Reading: loaded 5 lines from file 'data/name1.log'
2018/08/10 22:16:42 Saving: inserted 2 records into DB
2018/08/10 22:16:42 Total records in DB: 2
2018/08/10 22:16:42 Saving: inserted 5 records into DB
2018/08/10 22:16:42 Total records in DB: 7
```

If you open the files 'data/name1.log' and 'data/name2.log' in a text editor and duplicate one or more lines in them, you will see something like the following in the Terminal:

```shell
# Terminal:
...
2018/08/10 22:17:58 Reading: loaded 1 lines from file 'data/name1.log'
2018/08/10 22:17:59 Saving: inserted 1 records into DB
2018/08/10 22:17:59 Total records in DB: 8
2018/08/10 22:18:05 Reading: loaded 2 lines from file 'data/name1.log'
2018/08/10 22:18:05 Saving: inserted 2 records into DB
2018/08/10 22:18:05 Total records in DB: 10
2018/08/10 22:18:09 Reading: loaded 3 lines from file 'data/name2.log'
2018/08/10 22:18:09 Saving: inserted 3 records into DB
2018/08/10 22:18:09 Total records in DB: 13
2018/08/10 22:18:17 Reading: loaded 2 lines from file 'data/name2.log'
2018/08/10 22:18:17 Saving: inserted 2 records into DB
2018/08/10 22:18:17 Total records in DB: 15
```

## Graceful shutdown
If stop the application, you will see something like the following:

```shell
# Terminal:
2018/08/10 22:19:47 Signal: interrupt
2018/08/10 22:19:47 Terminating the application in 3 s
2018/08/10 22:19:48 Terminating the application in 2 s
2018/08/10 22:19:49 Terminating the application in 1 s
```
