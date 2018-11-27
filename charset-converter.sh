#!/bin/bash

if [ -z "$1" ]; then
	echo 1=$1
	exit
fi
if [ -z "$2" ]; then
	echo 2=$2
	exit
fi


find "$1" -type f -name "*.htm"|while read fname ;

do

   echo $fname

   echo " ---> $2/$fname"
   continue

   echo below continue

   #iconv -f GB18030 -t UTF-8 $fname > ${fname}.utf8

   #mv $fname ${fname}.18030

   #mv ${fname}.utf8 $fname

done
