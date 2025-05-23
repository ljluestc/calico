include ../metadata.mk

PACKAGE_NAME = github.com/projectcalico/calico/typha

# Name of the images.
# e.g., <registry>/<name>:<tag>
TYPHA_IMAGE    ?=typha
BUILD_IMAGES   ?=$(TYPHA_IMAGE)

###############################################################################
# Download and include ../lib.Makefile
#   Additions to EXTRA_DOCKER_ARGS need to happen before the include since
#   that variable is evaluated when we declare DOCKER_RUN and siblings.
###############################################################################
include ../lib.Makefile

###############################################################################

# All Typha go files.
SRC_FILES:=$(shell find . $(foreach dir,$(NON_TYPHA_DIRS),-path ./$(dir) -prune -o) -type f -name '*.go' -print)

# Touchfile created when we make the typha image.
TYPHA_CONTAINER_CREATED=.calico_typha.created-$(ARCH)
TYPHA_CONTAINER_FIPS_CREATED=.calico_typha.created-$(ARCH)-fips
FIPS ?= false

ifeq ($(FIPS),true)
TYPHA_CONTAINER_MARKER=$(TYPHA_CONTAINER_FIPS_CREATED)
VALIDARCHES=amd64
BINDIR=./bin/$(ARCH)-fips
else
TYPHA_CONTAINER_MARKER=$(TYPHA_CONTAINER_CREATED)
BINDIR=./bin
endif

.PHONY: clean
clean:
	# Clean .created files which indicate images / releases have been built.
	find . -name '.*.created*' -type f -delete
	find . -name '.*.published*' -type f -delete
	rm -rf bin \
		docker-image/bin \
		build \
		report/*.xml \
		release-notes-* \
		.calico_typha.created-*
	find . -name "*.coverprofile" -type f -delete
	find . -name "coverage.xml" -type f -delete
	find . -name ".coverage" -type f -delete
	find . -name "*.pyc" -type f -delete
	-docker image rm -f $$(docker images $(TYPHA_IMAGE) -a -q)


###############################################################################
# Building the binary
###############################################################################
build: $(BINDIR)/calico-typha-$(ARCH)
build-all: $(addprefix sub-build-,$(VALIDARCHES))
sub-build-%:
	$(MAKE) build ARCH=$*

$(BINDIR)/wrapper-$(ARCH): $(SRC_FILES) $(LOCAL_BUILD_DEP)
ifeq ($(FIPS),true)
	$(call build_cgo_boring_binary, ./cmd/wrapper, $@)
else
	$(call build_binary, ./cmd/wrapper, $@)
endif

$(BINDIR)/calico-typha-$(ARCH): $(SRC_FILES)
ifeq ($(FIPS),true)
	$(call build_cgo_boring_binary, ./cmd/calico-typha, $@)
else
	$(call build_binary, ./cmd/calico-typha, $@)
endif

$(BINDIR)/typha-client-$(ARCH): $(SRC_FILES)
	@echo Building typha client...
ifeq ($(FIPS),true)
	$(call build_cgo_boring_binary, ./cmd/typha-client, $@)
else
	$(call build_binary, ./cmd/typha-client, $@)
endif

###############################################################################
# Building the image
###############################################################################
# Build the calico/typha docker image, which contains only typha.
.PHONY: $(TYPHA_IMAGE) $(TYPHA_IMAGE)-$(ARCH)
image: $(BUILD_IMAGES)

# Build the image for the target architecture
.PHONY: image-all
image-all: $(addprefix sub-image-,$(VALIDARCHES)) sub-image-fips-amd64
sub-image-%:
	$(MAKE) image ARCH=$*
sub-image-fips-%:
	$(MAKE) image FIPS=true ARCH=$*

# Build the calico/typha docker image, which contains only Typha.
.PHONY: image $(TYPHA_IMAGE)
TYPHA_IMAGE_WITH_TAG=$(TYPHA_IMAGE):latest-$(ARCH)
TYPHA_IMAGE_ID=$(shell docker images -q $(TYPHA_IMAGE_WITH_TAG))

# This target actually builds the image if any of its dependencies have changed or if it's
# missing from docker.
$(TYPHA_CONTAINER_CREATED): $(BINDIR)/calico-typha-$(ARCH) \
                            LICENSE \
                            docker-image/Dockerfile \
                            docker-image/typha.cfg \
                            $(shell test "$(TYPHA_IMAGE_ID)" || echo force-rebuild)
	rm -rf docker-image/bin
	mkdir -p docker-image/bin
	cp $(BINDIR)/calico-typha-$(ARCH) docker-image/bin/
	cp LICENSE docker-image/
	$(DOCKER_BUILD) -t $(TYPHA_IMAGE):latest-$(ARCH) -f docker-image/Dockerfile docker-image
	$(MAKE) retag-build-images-with-registries VALIDARCHES=$(ARCH) IMAGETAG=latest
	touch $(TYPHA_CONTAINER_CREATED)

$(TYPHA_CONTAINER_FIPS_CREATED): $(BINDIR)/calico-typha-$(ARCH) \
                            LICENSE \
                            docker-image/Dockerfile \
                            docker-image/typha.cfg \
                            $(shell test "$(TYPHA_IMAGE_ID)" || echo force-rebuild)
	rm -rf docker-image/bin
	mkdir -p docker-image/bin
	cp $(BINDIR)/calico-typha-$(ARCH) docker-image/bin/
	cp LICENSE docker-image/
	$(DOCKER_BUILD) -t $(TYPHA_IMAGE):latest-fips-$(ARCH) -f docker-image/Dockerfile docker-image
	$(MAKE) retag-build-images-with-registries VALIDARCHES=$(ARCH) IMAGETAG=latest-fips LATEST_IMAGE_TAG=latest-fips
	touch $(TYPHA_CONTAINER_FIPS_CREATED)


# This target retags the already built image, deferring the build to the target above.
$(TYPHA_IMAGE): $(TYPHA_CONTAINER_MARKER)

###############################################################################
# Unit Tests
###############################################################################
.PHONY: ut
ut combined.coverprofile: $(SRC_FILES)
	@echo Running Go UTs.
	$(DOCKER_RUN) $(CALICO_BUILD) ./utils/run-coverage

###############################################################################
# CI/CD
###############################################################################
.PHONY: cd ci version
version: image
	docker run --rm $(TYPHA_IMAGE):latest-$(ARCH) calico-typha --version

ci: mod-download image-all version static-checks ut
ifeq (,$(filter k8sfv-test, $(EXCEPT)))
	@$(MAKE) k8sfv-test
endif

## Deploys images to registry
cd: image-all cd-common

fv: k8sfv-test

k8sfv-test: image
	$(MAKE) -C ../felix JUST_A_MINUTE=true USE_TYPHA=true FV_TYPHAIMAGE=$(TYPHA_IMAGE):latest-$(ARCH) TYPHA_VERSION=latest k8sfv-test

st:
	@echo "No STs available."

###############################################################################
# Release
###############################################################################
## Produces a clean build of release artifacts at the specified version.
release-build: .release-$(VERSION).created
.release-$(VERSION).created:
	$(MAKE) clean image-all RELEASE=true
	$(MAKE) retag-build-images-with-registries RELEASE=true IMAGETAG=$(VERSION)
	$(MAKE) retag-build-images-with-registries RELEASE=true IMAGETAG=latest
	$(MAKE) FIPS=true retag-build-images-with-registries RELEASE=true IMAGETAG=$(VERSION)-fips LATEST_IMAGE_TAG=latest-fips
	$(MAKE) FIPS=true retag-build-images-with-registries RELEASE=true IMAGETAG=latest-fips LATEST_IMAGE_TAG=latest-fips
	touch $@

## Pushes a github release and release artifacts produced by `make release-build`.
release-publish: release-prereqs .release-$(VERSION).published
.release-$(VERSION).published:
	$(MAKE) push-images-to-registries push-manifests IMAGETAG=$(VERSION) RELEASE=$(RELEASE) CONFIRM=$(CONFIRM)
	$(MAKE) FIPS=true push-images-to-registries push-manifests IMAGETAG=$(VERSION)-fips RELEASE=$(RELEASE) CONFIRM=$(CONFIRM)
	touch $@

# WARNING: Only run this target if this release is the latest stable release. Do NOT
# run this target for alpha / beta / release candidate builds, or patches to earlier Calico versions.
## Pushes `latest` release images. WARNING: Only run this for latest stable releases.
release-publish-latest: release-prereqs
	# Check latest versions match.
	if ! docker run $(TYPHA_IMAGE):latest-$(ARCH) calico-typha --version | grep '$(VERSION)'; then echo "Reported version:" `docker run $(TYPHA_IMAGE):latest-$(ARCH) calico-typha --version` "\nExpected version: $(VERSION)"; false; else echo "\nVersion check passed\n"; fi
	if ! docker run quay.io/$(TYPHA_IMAGE):latest-$(ARCH) calico-typha --version | grep '$(VERSION)'; then echo "Reported version:" `docker run quay.io/$(TYPHA_IMAGE):latest-$(ARCH) calico-typha --version` "\nExpected version: $(VERSION)"; false; else echo "\nVersion check passed\n"; fi

	$(MAKE) push-images-to-registries push-manifests IMAGETAG=latest RELEASE=$(RELEASE) CONFIRM=$(CONFIRM)

###############################################################################
# Developer helper scripts (not used by build or test)
###############################################################################
.PHONY: ut-no-cover
ut-no-cover: $(SRC_FILES)
	@echo Running Go UTs without coverage.
	$(DOCKER_RUN) $(CALICO_BUILD) ginkgo -r

.PHONY: ut-watch
ut-watch: $(SRC_FILES)
	@echo Watching go UTs for changes...
	$(DOCKER_RUN) $(CALICO_BUILD) ginkgo watch -r

# Launch a browser with Go coverage stats for the whole project.
.PHONY: cover-browser
cover-browser: combined.coverprofile
	go tool cover -html="combined.coverprofile"

.PHONY: cover-report
cover-report: combined.coverprofile
	# Print the coverage.  We use sed to remove the verbose prefix and trim down
	# the whitespace.
	@echo
	@echo ======== All coverage =========
	@echo
	@$(DOCKER_RUN) $(CALICO_BUILD) sh -c 'go tool cover -func combined.coverprofile | \
				   sed 's=$(PACKAGE_NAME)/==' | \
				   column -t'
	@echo
	@echo ======== Missing coverage only =========
	@echo
	@$(DOCKER_RUN) $(CALICO_BUILD) sh -c "go tool cover -func combined.coverprofile | \
				   sed 's=$(PACKAGE_NAME)/==' | \
				   column -t | \
				   grep -v '100\.0%'"

# Install or update the tools used by the build
.PHONY: update-tools
update-tools:
	go get -u github.com/onsi/ginkgo/ginkgo
