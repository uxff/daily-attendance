# daily-attendance
a Daily Attendance app writen by go, use beego and bootstrap frame work. support wechat OA accounts.

# Requirement

```
go get -u github.com/mattn/go-sqlite3
go get -u github.com/beego/i18n
go get -u github.com/mattn/go-runewidth
```

# How to Use

```
$ git clone git@github.com:uxff/daily-attendance.git
$ cd daily-attendance
#
# need node and npm
$ npm install -g bower
$ bower install
#
# build
$ go build
#
# you need start mysql service, and config mysql in:
$ vim conf/app.conf
# add line:
 datasource=root:password@tcp(127.0.0.1:3306)/beegoauth?charset=utf8

```

# How to Run

```
$ ./daily-attendance
```

# Preview
![](https://raw.githubusercontent.com/uxff/daily-attendance/master/20181127073913.png)
![](https://raw.githubusercontent.com/uxff/daily-attendance/master/20181127074015.png)


# Describe and todo implement

- wechat login
- email verify
- payment
- crontab
- wx oa api



