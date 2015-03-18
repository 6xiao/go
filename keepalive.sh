LOG=~
PROC=ls
FLAG='
-a
-l
'
NRPROC=`ps ax | grep -v grep | grep -w $PROC | grep -w "$FLAG" | wc -l`
if [ $NRPROC -lt 1 ]
then
echo $(date +%Y-%m-%d) $(date +%H:%M:%S) $PROC >> $LOG/restart.log
$PROC $FLAG >> $LOG/$(basename $PROC).$(date +%Y%m%d.%H%M%S.log) 2>&1 &
fi
