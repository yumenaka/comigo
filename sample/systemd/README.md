# Comigo Systemd Service

Simple systemd service configuration for running Comigo as a system service.

## Quick Setup

```bash
# 1. Copy the service file
sudo cp comigo.service /etc/systemd/system/

# 2. Edit configuration (adjust paths)
sudo nano /etc/systemd/system/comigo.service

# 3. Reload systemd
sudo systemctl daemon-reload

# 4. Enable and start
sudo systemctl enable comigo
sudo systemctl start comigo
```

## Configuration

Edit `/etc/systemd/system/comigo.service` and modify:

- `ExecStart`: Path to comigo binary and your book library
- `Environment`: Environment variables (port, language, etc.)

## Commands

| Command | Description |
|---------|-------------|
| `systemctl start comigo` | Start service |
| `systemctl stop comigo` | Stop service |
| `systemctl restart comigo` | Restart service |
| `systemctl status comigo` | Check status |
| `journalctl -u comigo -f` | View logs |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `COMIGO_PORT` | Server port | 1234 |
| `COMIGO_USERNAME` | Login username | - |
| `COMIGO_PASSWORD` | Login password | - |
