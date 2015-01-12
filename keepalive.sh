P=/home/user/process
NRP=`ps ax | grep -v grep | grep -w $P | wc -l`
if [ $NRP -lt 1 ]
then
echo $(date +%Y-%m-%d) $(date +%H:%M:%S) $P >> ~/restart.log
$P \
-flag="" \
>> /home/user/process/$(date +%Y%m%d.%H%M%S.log) 2>&1 &
fi
