# SSH Tunnel Bridger Utility

This utility allows you to establish SSH tunnels to containers running in **Docker** or **Kubernetes** 
environments, providing seamless access to remote services via local ports. 
The configuration is YAML-based, enabling clear, flexible, and reusable setups for different environments.

# Table of Contents

1. [Getting Started](#getting-started)
   - [Preparing SSH Access](#preparing-ssh-access)
   - [Write a Configuration File](#write-a-configuration-file)
   - [Build and Run the Bridger Container](#build-and-run-the-bridger-container)
2. [Examples](#examples)
   - [Tunnel to a Kubernetes Service](#example-1-tunnel-to-a-kubernetes-service)
   - [Tunnel Multiple Docker Services](#example-2-tunnel-multiple-docker-services)
3. [Logs and Debugging](#logs-and-debugging)
4. [Conclusion](#conclusion)

## Features

- **Docker and Kubernetes Support**: Tunnel to services running in both Docker containers and Kubernetes pods.
- **Environment-Specific Configurations**: Organize connections by environment (e.g., `local`, `uat`, `demo`).
- **Port Mapping**: Expose remote container ports locally via SSH.
- **Simplified CLI**: Specify configurations and mappings with an intuitive command-line interface.
- **Customizable Setup**: Define and manage multiple components in a single config file.

---

## Getting Started
[Back to Table of Contents](#table-of-contents)

### Prerequisites

- **Docker** installed on your local machine.
- SSH access configured for each environment in your config file, including a valid private key (`id_rsa`) and known hosts file.

---

### Preparing SSH Access
[Back to Table of Contents](#table-of-contents)

To use the SSH Tunnel Bridger, you need SSH access configured for each environment. This involves generating an SSH key and installing it on the remote environment.

#### Step 1: Generate an SSH Key

If you don't already have an SSH key, you can generate one using OpenSSH:

```bash
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```

- **Options**:
  - `-t rsa`: Specifies the RSA key type.
  - `-b 4096`: Creates a 4096-bit key for increased security.
  - `-C "your_email@example.com"`: Adds a label to the key for identification.

- **Follow the prompts**:
  - Save the key in the default location (`~/.ssh/id_rsa`) or specify a custom path.
  - Optionally, set a passphrase for additional security.

#### Step 2: Install the Key on a Remote Host

Copy the public key to the remote environment using the `ssh-copy-id` command:

```bash
ssh-copy-id username@remote-host
```

- Replace `username` with your SSH user and `remote-host` with the environment's hostname or IP address.
- The command appends your public key to the `~/.ssh/authorized_keys` file on the remote host.

#### Manual Installation (Alternative)

If `ssh-copy-id` isn't available, manually install the key:

1. Display your public key:
   ```bash
   cat ~/.ssh/id_rsa.pub
   ```

2. Log in to the remote host:
   ```bash
   ssh username@remote-host
   ```

3. Add the public key to the `~/.ssh/authorized_keys` file:
   ```bash
   echo "your-public-key" >> ~/.ssh/authorized_keys
   ```

4. Set the correct permissions:
   ```bash
   chmod 600 ~/.ssh/authorized_keys
   ```

#### Step 3: Verify the Connection

Test the SSH connection to ensure the key is working:

```bash
ssh username@remote-host
```

If successful, you should connect without being prompted for a password. You're now ready to use the SSH Tunnel Bridger!

### Step 1: Write a Configuration File
[Back to Table of Contents](#table-of-contents)

Create a YAML configuration file (`config.yaml`) to define environments and their services. Below is an example:

```yaml
- name: local
  hostname: michele@192.168.1.102
  components:
    - type: kubernetes
      name: nginx
      service: nginx-ingress-nginx-controller
      port: 80
      bridge-port: 10001
    - type: kubernetes
      name: whoami
      service: whoami
      port: 80
      bridge-port: 10002
- name: demo
  hostname: ubuntu@ukdemo.elerianai.com
  components:
    - type: docker
      name: ses
      service: aws-ses-middleware
      port: 8000
    - type: docker
      name: im
      service: intelligence-model-dashboard
      port: 8080
```

### Configuration Parameters

| Key          | Description                                                                 |
|--------------|-----------------------------------------------------------------------------|
| `name`       | Alias for the environment (e.g., `local`, `uat`).                           |
| `hostname`   | SSH user and host for the environment.                                      |
| `components` | List of services in the environment.                                        |
| `type`       | Type of service: `docker` or `kubernetes`.                                  |
| `name`       | (Optional) Human-readable name for the service.                            |
| `service`    | The Docker container name or Kubernetes service name to connect to.         |
| `port`       | The port exposed by the service.                                            |
| `bridge-port`| (Optional) Default local port mapping for the service.                     |

---

### Step 2: Build and Run the Bridger Container

Use the following Docker Compose configuration to run the SSH Tunnel Bridger:

```yaml
services:
  bridger:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 2345:2345
      - 1234:1234
    command:
      - bridger
      - /app/config.yaml
      - local:whoami:2345
    volumes:
      - ./config.yaml:/app/config.yaml
      - /home/michele/.ssh/known_hosts:/root/.ssh/known_hosts
      - /home/michele/.ssh/id_rsa:/root/.ssh/id_rsa
```

#### Explanation:
- **Command Structure**: 
  ```
  bridger <config file path> <service alias>:<environment>:<local port> ...
  ```
  Example:
  ```
  bridger /app/config.yaml local:whoami:2345
  ```
  This maps the `whoami` service in the `local` environment to `localhost:2345`.

- **Volumes**:
  - Map the `config.yaml` file into the container at `/app/config.yaml`.
  - Provide SSH credentials via mounted files.

---

### Examples

#### Example 1: Tunnel to a Kubernetes Service
YAML Config:
```yaml
- name: local
  hostname: michele@192.168.1.102
  components:
    - type: kubernetes
      name: nginx
      service: nginx-ingress-nginx-controller
      port: 80
      bridge-port: 10001
```

Command:
```bash
bridger /app/config.yaml local:nginx:8080
```
This maps the `nginx` service in `local` to `localhost:8080`.

#### Example 2: Tunnel Multiple Docker Services
YAML Config:
```yaml
- name: uat
  hostname: ubuntu@ukuat.elerianai.com
  components:
    - type: docker
      name: calm
      service: compose-calm-1
      port: 8000
      bridge-port: 10003
    - type: docker
      name: db
      service: compose-postgres-1
      port: 5432
```

Command:
```bash
bridger /app/config.yaml uat:calm:9000 uat:db:7000
```
This maps:
- `calm` service to `localhost:9000`.
- `db` service to `localhost:7000`.

---

### Logs and Debugging

The utility outputs detailed logs for each tunnel it establishes, helping you diagnose connection issues.

---

## Conclusion

This utility streamlines the process of accessing remote containerized services through SSH, making it easy to test and debug applications across Kubernetes and Docker environments. Customize your config file and start tunneling today!

