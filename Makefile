.PHONY: build install install-daemon

build:
	@rm -f build/mesh 2>/dev/null || true
	@go build -o build/mesh .

install: build
	@echo "📍 Installing binary..."
	@cp build/mesh /usr/local/bin/
	@chmod +x /usr/local/bin/mesh
	@echo "✅ Binary installed"

install-daemon: install
	@echo "⚙️  Installing systemd daemon..."
	@sudo tee /etc/systemd/system/mesh.service > /dev/null <<'EOF'
[Unit]
Description=MESH - Uptime Monitor Daemon
After=network.target

[Service]
Type=simple
User=$(shell whoami)
ExecStart=/usr/local/bin/mesh start
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF
	@sudo systemctl daemon-reload
	@sudo systemctl enable mesh.service
	@sudo systemctl start mesh.service
	@echo "✅ Daemon installed and started"
	@echo ""
	@echo "Status:"
	@sudo systemctl status mesh.service