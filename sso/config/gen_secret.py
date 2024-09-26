import os
import base64
import yaml


def generate_secret():
    secret = os.urandom(32)  # 256 bits
    return base64.urlsafe_b64encode(secret).decode('utf-8')


def update_yaml_file(secret, yaml_file='conf.yaml'):
    with open(yaml_file, 'r') as file:
        config = yaml.safe_load(file)

    if 'jwt' not in config:
        config['jwt'] = {}
    config['jwt']['secret'] = secret

    with open(yaml_file, 'w') as file:
        yaml.dump(config, file, default_flow_style=False)


if __name__ == "__main__":
    new_secret = generate_secret()
    update_yaml_file(new_secret)
    print(f"Updated JWT secret in conf.yaml")
