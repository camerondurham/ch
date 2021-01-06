all: $(DISTS)

BUILDS=\
  darwin-amd64  \
  linux-386     \
  linux-amd64   \
  linux-arm     \
  linux-arm64   \
  windows-386   \
  windows-amd64 \

dist:
	@mkdir -p dist

$(DISTS): OS = $(word 1,$(subst -, ,$*))
$(DISTS): ARCH = $(word 2,$(subst -, ,$*))

$(DISTS): dist/$(NAME)-%-$(VERSION).tgz: dist
	@echo "building: $@"
	@echo "OS = $(OS)"
	@echo "ARCH = $(ARCH)"
	@touch $@

cleanup:
	go mod tidy

get-tools:
	go get -t \
	github.com/spf13/cobra/cobra

todo:
	 git grep -EI "TODO|FIXME"

todos:
	 cp todos.txt todos.bkup.txt
	 git grep -EI "TODO|FIXME" > todos.txt
