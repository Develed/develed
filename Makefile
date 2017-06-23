VERTAG = $(shell grep "Version:" scripts/control | cut -c 10-)
BUILD = build/
IPKDIR = build/ipk/
DAEMONS = textd imaged dspd
TARGETS = proto bot $(DAEMONS)

.PHONY: $(TARGETS) release

all: $(TARGETS)

bot:
	@go build ./cmd/bot
dspd:
	@go build ./cmd/dspd
textd:
	@go build ./cmd/textd
imaged:
	@go build ./cmd/imaged
fake:
	@go build ./cmd/fake_req
proto:
	@protoc -I services/ services/services.proto --go_out=plugins=grpc:services

release: all
	@rm -rf $(BUILD)
	@mkdir -p $(IPKDIR)/usr/bin $(IPKDIR)/usr/share/develed $(IPKDIR)/etc/systemd/system
	@cp $(TARGETS) $(IPKDIR)/usr/bin/
	@cp config/sample.toml $(IPKDIR)/etc/develed.toml
	@cp -R resources/* $(IPKDIR)/usr/share/develed/
	@cp scripts/*.service $(IPKDIR)/etc/systemd/system/
	@cp scripts/control scripts/postinst scripts/preinst $(BUILD)
	@echo 2.0 > $(BUILD)/debian-binary
	@sed -i "s/:slack_bot_token:/$(SLACK_BOT_TOKEN)/g" $(IPKDIR)/etc/develed.toml
	@tar czf $(BUILD)/control.tar.gz -C $(BUILD) control postinst preinst
	@tar czf $(BUILD)/data.tar.gz -C $(IPKDIR) .
	@ar r $(BUILD)/develed_$(VERTAG).ipk $(BUILD)/control.tar.gz $(BUILD)/data.tar.gz $(BUILD)/debian-binary

test: all
	@killall $(TARGETS) >/dev/null 2>&1 || true
	@./dspd -config config/sample.toml -debug &
	@./imaged -config config/sample.toml &
	@./textd -config config/sample.toml &
	@sleep 0.5
	@./bot -config config/sample.toml -debug
	@killall $(DAEMONS) >/dev/null 2>&1

clean:
	rm -rf $(BUILD) $(TARGETS)
