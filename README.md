# daily-attendance
a Daily Attendance app writen by go, use beego and bootstrap framework. support wechat OA accounts.

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
#
# need node and npm and bower
# its ALREADY INSTALLED, no need do this
# $ npm install -g bower
# $ bower install
#
#
# build
$ go build
#
#
# you need start mysql service, and config mysql in:
$ vim conf/app.conf
# edit line:
 datasource=root:password@tcp(127.0.0.1:3306)/attendance?charset=utf8
# you must create database attendance firstly

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

# nginx 配置要求

需要nginx配置 proxy_set_header X-Host $host;proxy_set_header X-Scheme $scheme;

