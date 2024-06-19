#!/bin/bash
counter=1
while true
do
   ./main # 替换为你的二进制文件路径
   echo "第   $counter  次执行"
   ((counter++))
   sleep 8
done