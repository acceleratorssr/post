import os
import subprocess

output_dir = "../output"

if not os.path.exists(output_dir):
    os.makedirs(output_dir)

service = "../article"

print(f"Building -> {service} <-...")
os.chdir(service)

try:
    subprocess.run(["go", "build", "-o", output_dir, "."], check=True)
except subprocess.CalledProcessError as e:
    print(f"Error building {service}: {e}")
    exit(1)

os.chdir("..")

os.chdir("scripts")

print("sso services built successfully!")
