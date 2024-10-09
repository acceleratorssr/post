import os
import subprocess
import time

services_directory = "../output"
error_log = []
execution_order = ["sso.exe", "user.exe", "article.exe", "interactive.exe", "search.exe", "bff.exe"]

startupinfo = subprocess.STARTUPINFO()
startupinfo.dwFlags |= subprocess.STARTF_USESHOWWINDOW

for service in execution_order:
    exe_file = os.path.join(services_directory, service)

    if os.path.isfile(exe_file):
        print(f"Starting process for: {exe_file}")

        process = subprocess.Popen(['cmd.exe', '/c', 'start', '', exe_file], startupinfo=startupinfo)

        time.sleep(1)  # 等1s 没死就当它正常启动了

        print(f"{exe_file} is running in a new PowerShell window.")
    else:
        print(f"Executable not found: {exe_file}")
        error_log.append(f"Executable not found: {exe_file}")

if error_log:
    print(">! !-> Errors encountered during execution:")
    for err in error_log:
        print(err)
else:
    print("All processes started successfully.")
