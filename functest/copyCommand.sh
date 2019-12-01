#!/bin/bash
command=$1
basedir=/opt/golang/src/functest/chroottest
#copy commands
cmd=`which $command`
cmddir=${cmd%/*}
[ ! -d $basedir/$cmddir ] && mkdir $basedir/$cmddir -p 
cp $cmd $basedir/$cmddir

#copy libs
for line in $(ldd `which $command`  |grep lib |awk '{print $(NF-1)}')
do
	
	dir=${line%/*}
	[ ! -d $basedir/$dir ] && mkdir -p $basedir/$dir
	cp $line $basedir/$dir
done
