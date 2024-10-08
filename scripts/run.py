import os
import subprocess
import time

services_directory = "../output"
error_log = []
execution_order = ["sso.exe", "user.exe", "article.exe", "interactive.exe", "search.exe", "bff.exe"]

# 设置启动信息，以在新窗口中运行
startupinfo = subprocess.STARTUPINFO()
startupinfo.dwFlags |= subprocess.STARTF_USESHOWWINDOW

# 遍历执行顺序并启动进程
for service in execution_order:
    exe_file = os.path.join(services_directory, service)

    if os.path.isfile(exe_file):
        print(f"Starting process for: {exe_file}")

        # 启动进程并保持引用，使用 cmd /c start 来在新窗口中打开
        process = subprocess.Popen(['cmd.exe', '/c', 'start', '', exe_file], startupinfo=startupinfo)

        # 等待1秒
        time.sleep(1)

        # 这里不检查 process.poll()，因为我们已经在新窗口中启动了进程
        print(f"{exe_file} is running in a new PowerShell window.")
    else:
        print(f"Executable not found: {exe_file}")
        error_log.append(f"Executable not found: {exe_file}")

# 输出错误日志
if error_log:
    print(">! !-> Errors encountered during execution:")
    for err in error_log:
        print(err)
else:
    print("All processes started successfully.")
