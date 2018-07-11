VERTAG = $(shell grep "Version:" scripts/control | cut -c 10-)
BUILD = build/
IPKDIR = build/ipk/
DAEMONS = textd
TARGETS = bot $(DAEMONS)

# Verbose build?
ifeq ($(V),1)
Q :=
else
Q := @
endif

.PHONY: $(TARGETS) proto release

all: $(TARGETS)

bot:
	$Q go build ./cmd/bot
textd:
	$Q go build ./cmd/textd
fake:
	$Q go build ./cmd/fake_req
proto:
	$Q protoc -I services/ services/services.proto --go_out=plugins=grpc:services

release: all
	$Q rm -rf $(BUILD)
	$Q mkdir -p $(IPKDIR)/usr/bin $(IPKDIR)/usr/share/develed $(IPKDIR)/etc/systemd/system
	$Q cp $(TARGETS) $(IPKDIR)/usr/bin/
	$Q cp config/deploy.toml $(IPKDIR)/etc/develed.toml
	$Q cp -R resources/* $(IPKDIR)/usr/share/develed/
	$Q cp scripts/*.service $(IPKDIR)/etc/systemd/system/
	$Q cp scripts/control scripts/postinst scripts/preinst $(BUILD)
	$Q echo 2.0 > $(BUILD)/debian-binary
	$Q sed -i "s/:slack_bot_token:/$(SLACK_BOT_TOKEN)/g" $(IPKDIR)/etc/develed.toml
	$Q sed -i "s/:owm_token:/$(OWM_API_TOKEN)/g" $(IPKDIR)/etc/develed.toml
	$Q tar czf $(BUILD)/control.tar.gz -C $(BUILD) control postinst preinst
	$Q tar czf $(BUILD)/data.tar.gz -C $(IPKDIR) .
	$Q ar r $(BUILD)/develed_$(VERTAG).ipk $(BUILD)/control.tar.gz $(BUILD)/data.tar.gz $(BUILD)/debian-binary

test: all
	$Q killall $(TARGETS) >/dev/null 2>&1 || true
	$Q ./textd -config config/local.toml &
	$Q sleep 0.5
	$Q ./bot -config config/local.toml -debug
	$Q killall $(DAEMONS) >/dev/null 2>&1

clean:
	rm -rf $(BUILD) $(TARGETS)
