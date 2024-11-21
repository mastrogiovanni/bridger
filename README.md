```markdown
# SSH Tunnel Bridger Utility

This utility enables developers to establish SSH tunnels to containers running in **Docker** or **Kubernetes** environments, allowing seamless access to remote services via local ports. The configuration is YAML-based, providing clarity, flexibility, and reusability for different environments.

---

# Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Getting Started](#getting-started)
   - [Preparing SSH Access](#preparing-ssh-access)
   - [Write a Configuration File](#write-a-configuration-file)
   - [Build and Run the Bridger Container](#build-and-run-the-bridger-container)
4. [Examples](#examples)
   - [Tunnel to a Kubernetes Service](#example-1-tunnel-to-a-kubernetes-service)
   - [Tunnel Multiple Docker Services](#example-2-tunnel-multiple-docker-services)
5. [Logs and Debugging](#logs-and-debugging)
6. [Conclusion](#conclusion)

---

# Overview
[Back to Table of Contents](#table-of-contents)

The SSH Tunnel Bridger simplifies the process of accessing remote containerized services. By leveraging SSH tunneling, developers can map ports from Docker or Kubernetes services running in remote environments to their local machine. This ensures secure and efficient testing, debugging, and development workflows.

---

# Features
[Back to Table of Contents](#table-of-contents)

- **Docker and Kubernetes Support**: Tunnel to services running in both Docker containers and Kubernetes pods.
- **Environment-Specific Configurations**: Organize connections by environment (e.g., `dev`, `staging`, `prod`).
- **Port Mapping**: Expose remote container ports locally via SSH.
- **Simplified CLI**: Specify configurations and mappings with an intuitive command-line interface.
- **Customizable Setup**: Define and manage multiple components in a single configuration file.

---

# Getting Started
[Back to Table of Contents](#table-of-contents)

## Preparing SSH Access
[Back to Table of Contents](#table-of-contents)

To use the SSH Tunnel Bridger, you need SSH access configured for each environment. This involves generating an SSH key and installing it on the remote environment.

### Step 1: Generate an SSH Key

If you don't already have an SSH key, generate one using OpenSSH:

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

### Step 2: Install the Key on a Remote Host

Copy the public key to the remote environment using the `ssh-copy-id` command:

```bash
ssh-copy-id username@remote-host
```

- Replace `username` with your SSH user and `remote-host` with the environment's hostname or IP address.
- The command appends your public key to the `~/.ssh/authorized_keys` file on the remote host.

### Manual Installation (Alternative)

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

### Step 3: Verify the Connection

Test the SSH connection to ensure the key is working:

```bash
ssh username@remote-host
```

If successful, you should connect without being prompted for a password.

---

## Write a Configuration File
[Back to Table of Contents](#table-of-contents)

Create a YAML configuration file (`config.yaml`) to define environments and their services. Below is an example:

```yaml
- name: dev
  hostname: user@192.168.0.1
  components:
    - type: kubernetes
      name: api-service
      service: api-service
      port: 8080
      bridge-port: 10001
    - type: kubernetes
      name: web-service
      service: web-service
      port: 80
      bridge-port: 10002
- name: staging
  hostname: user@staging.example.com
  components:
    - type: docker
      name: app
      service: app-container
      port: 8000
    - type: docker
      name: db
      service: database-container
      port: 5432
```

### Configuration Parameters

| Key          | Description                                                                 |
|--------------|-----------------------------------------------------------------------------|
| `name`       | Alias for the environment (e.g., `dev`, `staging`, `prod`).                |
| `hostname`   | SSH user and host for the environment.                                      |
| `components` | List of services in the environment.                                        |
| `type`       | Type of service: `docker` or `kubernetes`.                                  |
| `name`       | (Optional) Human-readable name for the service.                            |
| `service`    | The Docker container name or Kubernetes service name to connect to.         |
| `port`       | The port exposed by the service.                                            |
| `bridge-port`| (Optional) Default local port mapping for the service.                     |

---

## Build and Run the Bridger Container
[Back to Table of Contents](#table-of-contents)

Use the following Docker Compose configuration to run the SSH Tunnel Bridger:

```yaml
services:
  bridger:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 9000:9000
      - 8000:8000
    command:
      - bridger
      - /app/config.yaml
      - dev:api-service:9000
    volumes:
      - ./config.yaml:/app/config.yaml
      - /home/user/.ssh/known_hosts:/root/.ssh/known_hosts
      - /home/user/.ssh/id_rsa:/root/.ssh/id_rsa
```

---

# Examples
[Back to Table of Contents](#table-of-contents)

## Example 1: Tunnel to a Kubernetes Service
[Back to Table of Contents](#table-of-contents)

YAML Config:
```yaml
- name: dev
  hostname: user@192.168.0.1
  components:
    - type: kubernetes
      name: api-service
      service: api-service
      port: 8080
      bridge-port: 10001
```

Command:
```bash
bridger /app/config.yaml dev:api-service:9000
```
This maps the `api-service` in the `dev` environment to `localhost:9000`.

## Example 2: Tunnel Multiple Docker Services
[Back to Table of Contents](#table-of-contents)

YAML Config:
```yaml
- name: staging
  hostname: user@staging.example.com
  components:
    - type: docker
      name: app
      service: app-container
      port: 8000
    - type: docker
      name: db
      service: database-container
      port: 5432
```

Command:
```bash
bridger /app/config.yaml staging:app:8080 staging:db:5433
```

---

# Logs and Debugging
[Back to Table of Contents](#table-of-contents)

The utility outputs detailed logs for each tunnel it establishes, helping you diagnose connection issues.

---

# Conclusion
[Back to Table of Contents](#table-of-contents)

The SSH Tunnel Bridger simplifies access to remote containerized services through SSH, making testing and debugging seamless across Kubernetes and Docker environments. Customize your configuration file and start tunneling today!
```