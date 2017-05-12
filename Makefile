VERTAG = $(shell grep "Version:" control | cut -c 10-)
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

release: all
	@rm -rf $(BUILD)
	@mkdir -p $(IPKDIR)/usr/bin $(IPKDIR)/usr/share/develed $(IPKDIR)/etc
	@cp dspd textd bot $(IPKDIR)/usr/bin/
	@cp config/sample.toml $(IPKDIR)/etc/develed.toml
	@cp -R cmd/textd/fonts/ $(IPKDIR)/usr/share/develed/fonts/
	@cp control $(BUILD)
	@echo 2.0 > $(BUILD)/debian-binary
	@tar czf $(BUILD)/control.tar.gz -C $(BUILD) control
	@tar czf $(BUILD)/data.tar.gz -C $(IPKDIR) .
	@ar r $(BUILD)/develed_$(VERTAG).ipk $(BUILD)/control.tar.gz $(BUILD)/data.tar.gz $(BUILD)/debian-binary

clean:
	rm -rf $(BUILD) dspd textd dspd
