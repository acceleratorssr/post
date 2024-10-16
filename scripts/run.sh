#!/bin/bash

servicesDirectory="../output"
errorLog=()
executionOrder=("sso" "user" "article" "interactive" "search" "bff" "recommend")

for service in "${executionOrder[@]}"; do
    exeFile="$servicesDirectory/$service"

    if [ -f "$exeFile" ]; then
        echo "Starting process for: $exeFile"

        gnome-terminal -- bash -c "$exeFile; exec bash" &> /dev/null &

        sleep 1

        # 检查进程是否仍在运行
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

if [ ${#errorLog[@]} -gt 0 ]; then
    echo ">! !-> Errors encountered during execution:"
    for err in "${errorLog[@]}"; do
        echo "$err"
    done
else
    echo "All processes started successfully."
fi
