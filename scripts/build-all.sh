#!/bin/bash

outputDir="../output"

# 检查输出目录是否存在，不存在则创建
if [ ! -d "$outputDir" ]; then
    mkdir -p "$outputDir"
fi

# 定义服务数组
services=("..\/article" "interactive" "search" "sso" "user" "bff")

# 遍历服务并构建
for service in "${services[@]}"; do
    echo "Building $service..."
    cd "$service" || { echo "Failed to enter directory $service"; exit 1; }

    # 执行构建命令，并直接检查结果
    if ! go build -o "$outputDir" .; then
        echo "Error building $service"
        exit 1
    fi

    cd - || exit 1
done

cd "scripts" || exit 1
echo "All services built successfully!"
