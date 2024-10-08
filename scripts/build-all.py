import os
import subprocess

output_dir = "../output"

# 检查输出目录是否存在，不存在则创建
if not os.path.exists(output_dir):
    os.makedirs(output_dir)

# 定义服务列表
services = ["../article", "interactive", "search", "sso", "user", "bff"]

# 遍历服务并构建
for service in services:
    print(f"Building {service}...")
    os.chdir(service)  # 切换到服务目录

    try:
        # 执行构建命令
        subprocess.run(["go", "build", "-o", output_dir, "."], check=True)
    except subprocess.CalledProcessError as e:
        print(f"Error building {service}: {e}")
        exit(1)

    os.chdir("..")  # 返回上一级目录

# 切换到脚本目录
os.chdir("scripts")

print("All services built successfully!")
