
GO       ?= go
GO_ENV   ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v2
NAME     ?= appbadger-cli
ARCH     ?= amd64
VERSION  ?= 0.1.0
MAIN     ?= *.go
BIN_DIR  ?= BIN
DEB_DIR  ?= DEB
BINARY   := $(BIN_DIR)/$(NAME)
BUILD_DATE ?= $(shell date -u +%Y-%m-%d)
BUILD_TIME ?= $(shell date -u +%H:%M:%S)

DIR_INSTALL ?= /usr/local/bin/
DIR_CONF ?= /etc/appbadger/
FILE_CONF ?= appbadger.json
BUILD_VERSION ?= 0.1.0

# Optimisations
BUILD_FLAGS ?= -trimpath
LDFLAGS     ?= -s -w

export fConfFiles:=$(fConfFiles)
export fControl:=$(fControl)
export fPreInst:=${fPreInst}
export fPostInst=$(fPostInst)
export fAppbadgerJson=$(fAppbadgerJson)

-include makefiles/*.mk

.PHONY: all build run clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO_ENV) $(GO) build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) \
			-X main.BUILD_VERSION=$(BUILD_VERSION) \
			-X main.BUILD_DATE=$(BUILD_DATE) \
			-X main.BUILD_TIME=$(BUILD_TIME)" \
		-o $(BINARY) $(MAIN)

clean:
	@rm -rf $(BIN_DIR) $(DEB_DIR)
	$(GO) clean

deb: clean build
	mkdir -p ${DEB_DIR}${DIR_INSTALL} \
		${DEB_DIR}${DIR_CONF} \
		${DEB_DIR}/DEBIAN

	cp -f ${BINARY} ${DEB_DIR}${DIR_INSTALL}${NAME} && \
	   chmod 440 ${DEB_DIR}${DIR_INSTALL}${NAME}

	printf "$${fAppbadgerJson}" > ${DEB_DIR}${DIR_CONF}${FILE_CONF}
	printf "$${fConfFiles}" > ${DEB_DIR}/DEBIAN/conffiles && /bin/chmod 0755 ${DEB_DIR}/DEBIAN/conffiles
	printf "$${fPreInst}" > ${DEB_DIR}/DEBIAN/preinst && /bin/chmod 0755 ${DEB_DIR}/DEBIAN/preinst
	printf "$${fPostInst}" > ${DEB_DIR}/DEBIAN/postinst && /bin/chmod 0755 ${DEB_DIR}/DEBIAN/postinst
	printf "$${fPreRm}" > ${DEB_DIR}/DEBIAN/prerm && /bin/chmod 0755 ${DEB_DIR}/DEBIAN/prerm
	printf "$${fPostRm}" > ${DEB_DIR}/DEBIAN/postrm && /bin/chmod 0755 ${DEB_DIR}/DEBIAN/postrm
	printf "$${fControl}" > ${DEB_DIR}/DEBIAN/control && /bin/chmod 0755 ${DEB_DIR}/DEBIAN/control

	fakeroot dpkg-deb -Zgzip --build ${DEB_DIR}/ ${DEB_DIR}/${NAME}_${VERSION}_${ARCH}.deb && \
	     fakeroot chown 1000:1000 ${DEB_DIR}/${NAME}_${VERSION}_${ARCH}.deb

	@if [ -d "${DEB_DIR}/" ]; then rm -rf ${DEB_DIR}/DEBIAN ${DEB_DIR}/etc ${DEB_DIR}/usr; fi
