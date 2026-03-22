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
	@sed "s/YOUR_USER/$(shell whoami)/" mesh.service | sudo tee /etc/systemd/system/mesh.service > /dev/null
	@sudo systemctl daemon-reload
	@sudo systemctl enable mesh.service
	@sudo systemctl start mesh.service
	@echo "✅ Daemon installed and started"
	@echo ""
	@echo "Status:"
	@sudo systemctl status mesh.service