include ../metadata.mk

CALICO_DIR=$(shell git rev-parse --show-toplevel)/calico
GIT_HASH=$(shell git rev-parse --short=9 HEAD)
VERSIONS_FILE?=$(CALICO_DIR)/_data/versions.yml
JEKYLL_VERSION=pages
HP_VERSION=v0.2
DEV?=false
CONFIG=--config _config.yml
ifeq ($(DEV),true)
	CONFIG:=$(CONFIG),_config_dev.yml
endif

# Set DEV_NULL=true to enable the Null Converter which renders the docs site as markdown.
# This is useful for comparing changes to templates & includes.
ifeq ($(DEV_NULL),true)
	CONFIG:=$(CONFIG),_config_null.yml
endif

include ../metadata.mk
CALICO_BUILD?=calico/go-build:$(GO_BUILD_VER)
LOCAL_USER_ID?=$(shell id -u $$USER)
PACKAGE_NAME = github.com/projectcalico/calico/calico

# Determine whether there's a local yaml installed or use dockerized version.
# Note in order to install local (faster) yaml: "go get github.com/mikefarah/yq.v2"
YAML_CMD:=$(shell which yq.v2 || echo docker run --rm -i mikefarah/yq:2.4.2 yq)

# Local directories to ignore when running htmlproofer
HP_IGNORE_LOCAL_DIRS="/v1.5/,/v1.6/,/v2.0/,/v2.1/,/v2.2/,/v2.3/,/v2.4/,/v2.5/,/v2.6/,/v3.0/"

##############################################################################
# Version information used for cutting a release.
RELEASE_STREAM := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].title' | grep --only-matching --extended-regexp '(v[0-9]+\.[0-9]+)|master')

# Use := so that these V_ variables are computed only once per make run.
CALICO_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].title')
NODE_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.calico/node.version')
CTL_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.calicoctl.version')
CNI_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.calico/cni.version')
KUBE_CONTROLLERS_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.calico/kube-controllers.version')
POD2DAEMON_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.flexvol.version')
DIKASTES_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.calico/dikastes.version')
FLANNEL_MIGRATION_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.calico/flannel-migration-controller.version')
TYPHA_VER := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].components.typha.version')
CHART_RELEASE := $(shell cat $(VERSIONS_FILE) | $(YAML_CMD) read - '[0].chart.version')

###############################################################################
# Include ../lib.Makefile
#   Additions to EXTRA_DOCKER_ARGS need to happen before the include since
#   that variable is evaluated when we declare DOCKER_RUN and siblings.
###############################################################################
include ../lib.Makefile

##############################################################################

CONTAINERIZED_VALUES?=docker run --rm -v $(CURDIR):/calico -w /calico ruby:2.5

# Build values.yaml for all charts
.PHONY: values.yaml
_includes/charts/%/values.yaml: _plugins/values.rb _plugins/helm.rb _data/versions.yml
	$(CONTAINERIZED_VALUES) ruby ./hack/gen_values_yml.rb --registry $(REGISTRY) --chart $* > $@

# The following chunk of conditionals sets the Version of the helm chart.
# Note that helm requires strict semantic versioning, so we use v0.0 to represent 'master'.
ifdef RELEASE_CHART
# the presence of RELEASE_CHART indicates we're trying to cut an official chart release.
chartVersion:=$(CALICO_VER)
appVersion:=$(CALICO_VER)
else
# otherwise, it's a nightly build.
ifeq ($(RELEASE_STREAM), master)
# For master, helm requires semantic versioning, so use v0.0
chartVersion:=v0.0
appVersion:=$(CALICO_VER)-$(GIT_HASH)
else
chartVersion:=$(RELEASE_STREAM)
appVersion:=$(CALICO_VER)-$(GIT_HASH)
endif
endif

charts: chart/tigera-operator
chart/%: _includes/charts/%/values.yaml bin/helm3
	mkdir -p bin
	bin/helm3 package ./_includes/charts/$(@F) \
	--destination ./bin/ \
	--version $(chartVersion) \
	--app-version $(appVersion)

serve: bin/helm
	# We have to override JEKYLL_DOCKER_TAG which is usually set to 'pages'.
	# When set to 'pages', jekyll starts in safe mode which means it will not
	# load any plugins. Since we're no longer running in github-pages, but would
	# like to use a docker image that comes preloaded with all the github-pages plugins,
	# its ok to override this variable.
	docker run --rm -it \
	  -v $(CURDIR)/bin/helm:/usr/local/bin/helm:ro \
	  -v $(CURDIR):/srv/jekyll \
	  -e JEKYLL_DOCKER_TAG="" \
	  -e JEKYLL_UID=`id -u` \
	  -p 4000:4000 \
	  jekyll/jekyll:$(JEKYLL_VERSION) jekyll serve --incremental $(CONFIG)

.PHONY: build
_site build: bin/helm
	docker run --rm \
	-e JEKYLL_DOCKER_TAG="" \
	-e JEKYLL_UID=`id -u` \
	-v $(CURDIR)/bin/helm:/usr/local/bin/helm:ro \
	-v $(CURDIR):/srv/jekyll \
	-v $(VERSIONS_FILE):/srv/jekyll/_data/versions.yml \
	jekyll/jekyll:$(JEKYLL_VERSION) jekyll build --incremental $(CONFIG)

## Clean enough that a new release build will be clean
clean:
	rm -rf _output _site .jekyll-metadata pinned_versions.yaml _includes/charts/*/values.yaml

###############################################################################
# CI / test targets
###############################################################################

ci: htmlproofer kubeval helm-tests

htmlproofer: _site
	docker run -ti -e JEKYLL_UID=`id -u` --rm \
		-v $(CURDIR)/_site:/_site/ \
		quay.io/calico/htmlproofer:$(HP_VERSION) \
		/_site --assume-extension --check-html --empty-alt-ignore --file-ignore $(HP_IGNORE_LOCAL_DIRS) --internal_domains "projectcalico.docs.tigera.io" --disable_external --allow-hash-href

kubeval: _site
	# Run kubeval to check master manifests are valid Kubernetes resources.
	-docker run -v $(CURDIR):/calico --entrypoint /bin/sh garethr/kubeval:0.7.3 -c 'ok=true; for f in `find /calico/_site/master -name "*.yaml" |grep -v "\(config\|allow-istio-pilot\|30-policy\|istio-app-layer-policy\|istio-inject-configmap.*\|-cf\).yaml"`; do echo Running kubeval on $$f; /kubeval $$f || ok=false; done; $$ok' 1>stderr.out 2>&1

	# Filter out error loading schema for non-standard resources.
	# Filter out error reading empty secrets (which we use for e.g. etcd secrets and seem to work).
	-grep -v "Could not read schema from HTTP, response status is 404 Not Found" stderr.out | grep -v "invalid Secret" > filtered.out

	# Display the errors with context and fail if there were any.
	-rm stderr.out
	! grep -C3 -P "invalid|\t\*" filtered.out
	rm filtered.out

helm-tests: bin/helm values.yaml
	$(DOCKER_RUN) \
		-v $(CURDIR)/bin/helm3:/usr/local/bin/helm \
		$(CALICO_BUILD) \
		sh -c "$(GIT_CONFIG_SSH) ginkgo -cover -r ./helm-tests -chart-path=./_includes/$(RELEASE_STREAM)/charts/calico $(GINKGO_ARGS)"

###############################################################################
# Docs automation
###############################################################################

# URLs to ignore when checking external links.
HP_IGNORE_URLS="/docs.openshift.org/,/localhost/"

check_external_links: _site
	docker run -ti -e JEKYLL_UID=`id -u` --rm -v $(CURDIR)/_site:/_site/ quay.io/calico/htmlproofer:$(HP_VERSION) /_site --external_only --file-ignore $(HP_IGNORE_LOCAL_DIRS) --assume-extension --url-ignore $(HP_IGNORE_URLS) --internal_domains "docs.projectcalico.org"

strip_redirects:
	find \( -name '*.md' -o -name '*.html' \) -exec sed -i'' '/redirect_from:/d' '{}' \;

add_redirects_for_latest: strip_redirects
ifndef VERSION
	$(error VERSION is undefined - run using make add_redirects_for_latest VERSION=vX.Y)
endif
	# Check that the VERSION directory already exists
	@test -d $(VERSION)

	# Add the redirect line - look at .md files only and add "redirect_from: XYZ" on a new line after each "title:"
	find $(VERSION) \( -name '*.md' -o -name '*.html' \) -exec sed -i 's#^title:.*#&\nredirect_from: {}#' '{}' \;

	# Check the redirect_from lines and update the version to be "latest"
	find $(VERSION) \( -name '*.md' -o -name '*.html' \) -exec sed -i 's#^\(redirect_from: \)$(VERSION)#\1latest#' '{}' \;

	# Check the redirect_from lines and strip the .md from the URL
	find $(VERSION) \( -name '*.md' -o -name '*.html' \) -exec sed -i 's#^\(redirect_from:.*\)\.md#\1#' '{}' \;

update_canonical_urls:
	# Looks through all directories and replaces previous latest release version numbers in canonical URLs with new
	python release-scripts/update-canonical-urls.py

###############################################################################
# Release targets
###############################################################################

## Tags and builds a release from start to finish.
release: release-prereqs
	$(MAKE) RELEASE_CHART=true release-tag
	$(MAKE) RELEASE_CHART=true release-build
	$(MAKE) RELEASE_CHART=true release-verify

	@echo ""
	@echo "Release build complete. Next, push the release."
	@echo ""
	@echo "  make release-publish"
	@echo ""

## Produces a git tag for the release.
release-tag: release-prereqs
	git tag $(CALICO_VER)

## Produces a clean build of release artifacts at the specified version.
release-build: release-prereqs clean
	# Create the release archive.
	$(MAKE) release-archive

## Verifies the release artifacts produces by `make release-build` are correct.
release-verify: release-prereqs
	@echo "TODO: Implement release tar verification"

ifneq (,$(findstring $(RELEASE_STREAM),v3.5 v3.4 v3.3 v3.2 v3.1 v3.0 v2.6))
    # Found: this is an older release.
    REL_NOTES_PATH:=releases
else
    # Not found: this is a newer release.
    REL_NOTES_PATH:=release-notes
endif

UPLOAD_DIR?=$(OUTPUT_DIR)/upload
$(UPLOAD_DIR):
	mkdir -p $(UPLOAD_DIR)

# Define a multi-line string for the GitHub release body.
# We need to export it as an env var to properly format it.
# See here: https://stackoverflow.com/questions/649246/is-it-possible-to-create-a-multi-line-string-variable-in-a-makefile/5887751
define RELEASE_BODY
Release notes can be found at https://projectcalico.docs.tigera.io/archive/$(RELEASE_STREAM)/$(REL_NOTES_PATH)/

Attached to this release are the following artifacts:

- `release-v$(CALICO_VER).tgz`: docker images and kubernetes manifests.
- `calico-windows-v$(CALICO_VER).zip`: Calico for Windows.
- `tigera-operator-v$(CALICO_VER)-$(CHART_RELEASE).tgz`: Calico helm v3 chart.

endef
export RELEASE_BODY

## Pushes a github release and release artifacts produced by `make release-build`.
release-publish: release-prereqs $(UPLOAD_DIR)
	# Push the git tag.
	git push origin $(CALICO_VER)

	cp $(RELEASE_HELM_CHART) $(RELEASE_DIR).tgz $(RELEASE_WINDOWS_ZIP) $(UPLOAD_DIR)

	# Push binaries to GitHub release.
	# Requires ghr: https://github.com/tcnksm/ghr
	# Requires GITHUB_TOKEN environment variable set.
	ghr -u projectcalico -r calico \
		-b "$$RELEASE_BODY" \
		-n $(CALICO_VER) \
		$(CALICO_VER) $(UPLOAD_DIR)

	$(MAKE) helm-index
	@echo "Verify the GitHub release based on the pushed tag."
	@echo ""
	@echo "  https://github.com/projectcalico/calico/calico/releases/tag/$(CALICO_VER)"
	@echo ""

## Kicks semaphore job which syncs github released helm charts with helm index file
.PHONY: helm-index
helm-index:
	@echo "Triggering semaphore workflow to update helm index."
	SEMAPHORE_PROJECT_ID=30f84ab3-1ea9-4fb0-8459-e877491f3dea SEMAPHORE_WORKFLOW_BRANCH=master SEMAPHORE_WORKFLOW_FILE=../releases/calico/helmindex/update_helm.yml $(MAKE) semaphore-run-workflow

## Generates release notes for the given version.
.PHONY: release-notes
release-notes: #release-prereqs
	VERSION=$(CALICO_VER) GITHUB_TOKEN=$(GITHUB_TOKEN) python2 ./release-scripts/generate-release-notes.py

update-authors:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN must be set)
endif
	@echo "# Calico authors" > AUTHORS.md
	@echo "" >> AUTHORS.md
	@echo "This file is auto-generated based on contribution records reported" >> AUTHORS.md
	@echo "by GitHub for the core repositories within the projectcalico/ organization. It is ordered alphabetically." >> AUTHORS.md
	@echo "" >> AUTHORS.md
	@docker run -ti --rm -v $(CURDIR):/code -e GITHUB_TOKEN=$(GITHUB_TOKEN) python:3 \
		bash -c 'pip install pygithub && /usr/local/bin/python /code/release-scripts/get-contributors.py >> /code/AUTHORS.md'

# release-prereqs checks that the environment is configured properly to create a release.
release-prereqs: charts
	@if [ $(CALICO_VER) != $(NODE_VER) ]; then \
		echo "Expected CALICO_VER $(CALICO_VER) to equal NODE_VER $(NODE_VER)"; \
		exit 1; fi
ifeq (, $(shell which ghr))
	$(error Unable to find `ghr` in PATH, run this: go get -u github.com/tcnksm/ghr)
endif

OUTPUT_DIR?=_output
RELEASE_DIR_NAME?=release-$(CALICO_VER)
RELEASE_DIR?=$(OUTPUT_DIR)/$(RELEASE_DIR_NAME)
RELEASE_DIR_K8S_MANIFESTS?=$(RELEASE_DIR)/k8s-manifests
RELEASE_DIR_IMAGES?=$(RELEASE_DIR)/images
RELEASE_DIR_BIN?=$(RELEASE_DIR)/bin
RELEASE_WINDOWS_ZIP=$(OUTPUT_DIR)/calico-windows-$(NODE_VER).zip
RELEASE_HELM_CHART=bin/tigera-operator-$(CALICO_VER)-$(CHART_RELEASE).tgz

# Determine where the manifests live. For older versions we used
# a different location, but we still need to package them up for patch
# releases.
DEFAULT_MANIFEST_SRC=./_site/manifests
OLD_VERSIONS := v3.0 v3.1 v3.2 v3.3 v3.4 v3.5 v3.6
ifneq ($(filter $(RELEASE_STREAM),$(OLD_VERSIONS)),)
DEFAULT_MANIFEST_SRC=./_site/$(RELEASE_STREAM)/getting-started/kubernetes/installation
endif
MANIFEST_SRC?=$(DEFAULT_MANIFEST_SRC)

$(RELEASE_WINDOWS_ZIP):
	wget https://github.com/projectcalico/calico/node/releases/download/$(NODE_VER)/calico-windows-$(NODE_VER).zip -P $(OUTPUT_DIR)

## Create an archive that contains a complete "Calico" release. This includes the release tarball (which bundles manifests, images, and binaries) and the Calico for Windows installation archive.
release-archive: release-prereqs $(RELEASE_DIR).tgz $(RELEASE_WINDOWS_ZIP)

$(RELEASE_DIR).tgz: $(RELEASE_DIR) $(RELEASE_DIR_K8S_MANIFESTS) $(RELEASE_DIR_IMAGES) $(RELEASE_DIR_BIN) $(RELEASE_DIR)/README
	tar -czvf $(RELEASE_DIR).tgz -C $(OUTPUT_DIR) $(RELEASE_DIR_NAME)

$(RELEASE_DIR_IMAGES): $(RELEASE_DIR_IMAGES)/calico-node.tar $(RELEASE_DIR_IMAGES)/calico-typha.tar $(RELEASE_DIR_IMAGES)/calico-cni.tar $(RELEASE_DIR_IMAGES)/calico-kube-controllers.tar $(RELEASE_DIR_IMAGES)/calico-pod2daemon-flexvol.tar $(RELEASE_DIR_IMAGES)/calico-dikastes.tar $(RELEASE_DIR_IMAGES)/calico-flannel-migration-controller.tar


$(RELEASE_DIR_BIN): $(RELEASE_DIR_BIN)/calicoctl $(RELEASE_DIR_BIN)/calicoctl-windows-amd64.exe $(RELEASE_DIR_BIN)/calicoctl-darwin-amd64

$(RELEASE_DIR)/README:
	@echo "This directory contains a complete release of Calico $(CALICO_VER)" >> $@
	@echo "Documentation for this release can be found at http://projectcalico.docs.tigera.io/$(RELEASE_STREAM)" >> $@
	@echo "" >> $@
	@echo "Docker images (under 'images'). Load them with 'docker load'" >> $@
	@echo "* The calico/node docker image  (version $(NODE_VERS))" >> $@
	@echo "* The calico/typha docker image  (version $(TYPHA_VER))" >> $@
	@echo "* The calico/cni docker image  (version $(CNI_VERS))" >> $@
	@echo "* The calico/kube-controllers docker image (version $(KUBE_CONTROLLERS_VER))" >> $@
	@echo "* The calico/dikastes docker image (version $(DIKASTES_VER))" >> $@
	@echo "* The calico/pod2daemon-flexvol docker image (version $(POD2DAEMON_VER))" >> $@
	@echo "* The calico/flannel-migration-controller docker image (version $(FLANNEL_MIGRATION_VER))" >> $@
	@echo "" >> $@
	@echo "Binaries (for amd64) (under 'bin')" >> $@
	@echo "* The calicoctl binary (for Linux) (version $(CTL_VER))" >> $@
	@echo "* The calicoctl-windows-amd64.exe binary (for Windows) (version $(CTL_VER))" >> $@
	@echo "* The calicoctl-darwin-amd64 binary (for Mac) (version $(CTL_VER))" >> $@
	@echo "" >> $@
	@echo "Kubernetes manifests (under 'k8s-manifests directory')" >> $@

$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)

$(RELEASE_DIR_K8S_MANIFESTS):
	# Ensure that the docs site is generated
	rm -rf ../_site
	$(MAKE) _site

	# Find all the hosted manifests and copy them into the release dir. Use xargs to mkdir the destination directory structure before copying them.
	# -printf "%P\n" prints the file name and directory structure with the search dir stripped off
	find $(MANIFEST_SRC) -name  '*.yaml' -printf "%P\n" | \
	  xargs -I FILE sh -c \
	    'mkdir -p $(RELEASE_DIR_K8S_MANIFESTS)/`dirname FILE`;\
	    cp $(MANIFEST_SRC)/FILE $(RELEASE_DIR_K8S_MANIFESTS)/`dirname FILE`;'

$(RELEASE_DIR_IMAGES)/calico-node.tar:
	mkdir -p $(RELEASE_DIR_IMAGES)
	docker pull calico/node:$(NODE_VER)
	docker save --output $@ calico/node:$(NODE_VER)

$(RELEASE_DIR_IMAGES)/calico-typha.tar:
	mkdir -p $(RELEASE_DIR_IMAGES)
	docker pull calico/typha:$(TYPHA_VER)
	docker save --output $@ calico/typha:$(TYPHA_VER)

$(RELEASE_DIR_IMAGES)/calico-cni.tar:
	mkdir -p $(RELEASE_DIR_IMAGES)
	docker pull calico/cni:$(CNI_VER)
	docker save --output $@ calico/cni:$(CNI_VER)

$(RELEASE_DIR_IMAGES)/calico-kube-controllers.tar:
	mkdir -p $(RELEASE_DIR_IMAGES)
	docker pull calico/kube-controllers:$(KUBE_CONTROLLERS_VER)
	docker save --output $@ calico/kube-controllers:$(KUBE_CONTROLLERS_VER)

$(RELEASE_DIR_IMAGES)/calico-pod2daemon-flexvol.tar:
	mkdir -p $(RELEASE_DIR_IMAGES)
	docker pull calico/pod2daemon-flexvol:$(POD2DAEMON_VER)
	docker save --output $@ calico/pod2daemon-flexvol:$(POD2DAEMON_VER)

$(RELEASE_DIR_IMAGES)/calico-dikastes.tar:
	mkdir -p $(RELEASE_DIR_IMAGES)
	docker pull calico/dikastes:$(DIKASTES_VER)
	docker save --output $@ calico/dikastes:$(DIKASTES_VER)

$(RELEASE_DIR_IMAGES)/calico-flannel-migration-controller.tar:
	mkdir -p $(RELEASE_DIR_IMAGES)
	docker pull calico/flannel-migration-controller:$(FLANNEL_MIGRATION_VER)
	docker save --output $@ calico/flannel-migration-controller:$(FLANNEL_MIGRATION_VER)

$(RELEASE_DIR_BIN)/%:
	mkdir -p $(RELEASE_DIR_BIN)
	wget https://github.com/projectcalico/calico/calicoctl/releases/download/$(CTL_VER)/$(@F) -O $@
	chmod +x $@

###############################################################################
# Utilities
###############################################################################
# TODO: stop using bin/helm as an entrypoint in build scripts.
bin/helm: bin/helm3
	mkdir -p bin
	$(eval TMP := $(shell mktemp -d))
	wget -q https://get.helm.sh/helm-v2.16.3-linux-amd64.tar.gz -O $(TMP)/helm.tar.gz
	tar -zxvf $(TMP)/helm.tar.gz -C $(TMP)
	mv $(TMP)/linux-amd64/helm bin/helm

helm-deps: bin/helm3 bin/helm
bin/helm3:
	mkdir -p bin
	$(eval TMP := $(shell mktemp -d))
	wget -q https://get.helm.sh/helm-v3.3.1-linux-amd64.tar.gz -O $(TMP)/helm3.tar.gz
	tar -zxvf $(TMP)/helm3.tar.gz -C $(TMP)
	mv $(TMP)/linux-amd64/helm bin/helm3

.PHONY: values.yaml
values.yaml: _includes/charts/calico/values.yaml _includes/charts/tigera-operator/values.yaml
_includes/charts/%/values.yaml: _plugins/values.rb _plugins/helm.rb _data/versions.yml
	docker run --rm \
	  -v $(CURDIR):/calico \
	  -w /calico \
	  ruby:2.5 ruby ./hack/gen_values_yml.rb --chart $* > $@

DOCS_TEST_CONTAINER=projectcalico/release-test
.PHONY: release-test-image
release-test-image:
	cd release-scripts/tests && docker build -t $(DOCS_TEST_CONTAINER) . && cd -

.PHONY: release-test
release-test: release-test-image
	docker run --rm \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v $(CURDIR):/docs \
	-e RELEASE_STREAM=$(RELEASE_STREAM) \
	$(DOCS_TEST_CONTAINER) sh -c \
	"nosetests . -e "$(EXCLUDE_REGEX)" \
	-s -v --with-xunit \
	--xunit-file='/docs/nosetests.xml' \
	--with-timer $(EXTRA_NOSE_ARGS)"

API_GEN_REPO?=tmjd/gen-crd-api-reference-docs
API_GEN_BRANCH?=kb_v2
OPERATOR_VERSION?=master
OPERATOR_REPO?=tigera/operator
build-operator-reference:
	mkdir -p .go-pkg-cache && \
	   docker run --rm \
	   --net=host \
	   -v $(CURDIR):/go/src/$(PACKAGE_NAME):rw \
	   -v $(CURDIR)/.go-pkg-cache:/go/pkg:rw \
	   -e LOCAL_USER_ID=$(LOCAL_USER_ID) \
	   -w /go/src/$(PACKAGE_NAME) \
	   $(CALICO_BUILD) /bin/bash -c 'export GO111MODULE=on && rm -rf builder && mkdir builder && cd builder && \
	           git clone --depth=1 -b $(API_GEN_BRANCH) https://github.com/$(API_GEN_REPO) api-gen && cd api-gen && \
	           go mod edit -replace github.com/tigera/operator=github.com/$(OPERATOR_REPO)@$(OPERATOR_VERSION) && \
	           go mod download && go build && \
	           ./gen-crd-api-reference-docs -config /go/src/$(PACKAGE_NAME)/reference/installation/config.json \
	                   -api-dir github.com/tigera/operator/api -out-file /go/src/$(PACKAGE_NAME)/reference/installation/_api.html'
