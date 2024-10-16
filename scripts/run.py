import os
import subprocess
import time
import platform

services_directory = "../output"
error_log = []
execution_order = ["sso", "user", "article", "interactive", "search", "recommend", "bff"]

current_platform = platform.system()

for service in execution_order:
    exe_file = os.path.join(services_directory, service)

    if current_platform == "Windows":
        exe_file += ".exe"

    if os.path.isfile(exe_file):
        print(f"Starting process for: {exe_file}")

        if current_platform == "Windows":
            process = subprocess.Popen(['cmd.exe', '/c', 'start', '', exe_file])
            print(f"{exe_file} is running in a new Command Prompt window.")
        else:
            process = subprocess.Popen([exe_file])
            print(f"{exe_file} is running.")

        time.sleep(1)  # 等待1秒，假设它正常启动
    else:
        print(f"Executable not found: {exe_file}")
        error_log.append(f"Executable not found: {exe_file}")

# 打印错误日志
if error_log:
    print(">! !-> Errors encountered during execution:")
    for err in error_log:
        print(err)
else:
    print("All processes started successfully.")
