.PHONY: init
init:
ifeq ($(shell uname -s),Darwin)
	@grep -r -l ghcli-enance * .goreleaser.yml | xargs sed -i "" "s/ghcli-enance/$$(basename `git rev-parse --show-toplevel`)/"
else
	@grep -r -l ghcli-enance * .goreleaser.yml | xargs sed -i "s/ghcli-enance/$$(basename `git rev-parse --show-toplevel`)/"
endif
