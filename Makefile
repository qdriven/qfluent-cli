.PHONY: init
init:
ifeq ($(shell uname -s),Darwin)
	@grep -r -l qfluent-cli * .goreleaser.yml | xargs sed -i "" "s/qfluent-cli/$$(basename `git rev-parse --show-toplevel`)/"
else
	@grep -r -l qfluent-cli * .goreleaser.yml | xargs sed -i "s/qfluent-cli/$$(basename `git rev-parse --show-toplevel`)/"
endif
