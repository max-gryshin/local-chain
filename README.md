# local-chain

### Requirements
 To make it possible deploy project install:
 [docker](https://docs.docker.com/engine/install/ubuntu/)
 [task](https://taskfile.dev/#/installation)

Do not forget to add your user to docker group:
```bash
sudo usermod -aG docker $USER
```
To run docker without `sudo`. Then restart the computer.
### Getting started
 To run project locally:
 ```bash
 task up
 ```
 Then we need to run migrations for the database
 ```bash
 task migration-up
 ```

To stop project run
```bash
task down
```