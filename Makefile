VERTAG = $(shell grep "Version:" scripts/control | cut -c 10-)
BUILD = build/
IPKDIR = build/ipk/

.PHONY: bot textd dspd release

all: bot textd dspd

bot:
	@go build ./cmd/bot
dspd:
	@go build ./cmd/dspd
textd:
	@go build ./cmd/textd
fake:
	@go build ./cmd/fake_req

release: all
	@rm -rf $(BUILD)
	@mkdir -p $(IPKDIR)/usr/bin $(IPKDIR)/usr/share/develed $(IPKDIR)/etc/systemd/system
	@cp dspd textd bot $(IPKDIR)/usr/bin/
	@cp config/sample.toml $(IPKDIR)/etc/develed.toml
	@cp -R resources/* $(IPKDIR)/usr/share/develed/
	@cp scripts/*.service $(IPKDIR)/etc/systemd/system/
	@cp scripts/control scripts/postinst scripts/preinst $(BUILD)
	@echo 2.0 > $(BUILD)/debian-binary
	@sed -i "s/:slack_bot_token:/$(SLACK_BOT_TOKEN)/g" $(IPKDIR)/etc/develed.toml
	@tar czf $(BUILD)/control.tar.gz -C $(BUILD) control postinst preinst
	@tar czf $(BUILD)/data.tar.gz -C $(IPKDIR) .
	@ar r $(BUILD)/develed_$(VERTAG).ipk $(BUILD)/control.tar.gz $(BUILD)/data.tar.gz $(BUILD)/debian-binary

clean:
	rm -rf $(BUILD) dspd textd dspd
