# Inferno 🔥

**Inferno** is the control plane for a distributed tunnel management system, orchestrating [Spark](https://github.com/vayzur/spark) agents and managing tunnels powered by [Xray-core](https://github.com/XTLS/Xray-core).

---

## 🛠 Development
### Requirements
- Go 1.23+
- etcd 3.6+ installed locally for testing

### Build
```bash
git clone https://github.com/vayzur/inferno.git
cd inferno

go build -o inferno cmd/inferno/main.go
```
