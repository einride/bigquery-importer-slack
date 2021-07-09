ko_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
ko_version := 0.6.0
ko := $(ko_cwd)/$(ko_version)/ko

ifeq ($(shell uname),Linux)
ko_archive_url := https://github.com/google/ko/releases/download/v$(ko_version)/ko_$(ko_version)_Linux_x86_64.tar.gz
else ifeq ($(shell uname),Darwin)
ko_archive_url := https://github.com/google/ko/releases/download/v$(ko_version)/ko_$(ko_version)_Darwin_x86_64.tar.gz
else
$(error unsupported OS: $(shell uname))
endif

$(ko):
	$(info building ko...)
	@mkdir -p $(dir $@)
	@curl -sSL $(ko_archive_url) -o - | tar -xz --directory $(dir $@)
	@chmod +x $@
	@touch $@
