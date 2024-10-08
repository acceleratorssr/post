#!/bin/bash

servicesDirectory="../output"
errorLog=()
executionOrder=("sso.exe" "user.exe" "article.exe" "interactive.exe" "search.exe" "bff.exe")

# 遍历执行顺序并启动进程
for service in "${executionOrder[@]}"; do
    exeFile="$servicesDirectory/$service"

    if [ -f "$exeFile" ]; then
        echo "Starting process for: $exeFile"

        # 使用 gnome-terminal 启动新窗口运行进程
        gnome-terminal -- bash -c "$exeFile; exec bash" &> /dev/null &

        # 等待1秒
        sleep 1

        # 检查进程是否仍在运行（使用 pgrep）
        processId=$!
        if ps -p $processId > /dev/null; then
            echo "$exeFile is running successfully."
        else
            errorMsg="$exeFile exited unexpectedly."
            echo "$errorMsg"
            errorLog+=("$errorMsg")
        fi
    else
        echo "Executable not found: $exeFile"
        errorLog+=("Executable not found: $exeFile")
    fi
done

# 输出错误日志
if [ ${#errorLog[@]} -gt 0 ]; then
    echo ">! !-> Errors encountered during execution:"
    for err in "${errorLog[@]}"; do
        echo "$err"
    done
else
    echo "All processes started successfully."
fi
