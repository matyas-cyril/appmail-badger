define fControl
Package: ${NAME}
Version: ${BUILD_VERSION}
Maintainer: matyas-cyril
Description: appbadger command line interface
Section: base
Priority: optional
Architecture: amd64
Installed-Size: $(du -s ${DEB_DIR} 2> /dev/null | cut -f1)
Depends: man

endef

define fPreInst
#!/bin/bash
mkdir -p ${DIR_INSTALL} ${DIR_CONF}

endef

define fPostInst
#!/bin/bash
chown root:root ${DIR_INSTALL}${APP} && /bin/chmod 0440 ${DIR_INSTALL}${APP}
chown root:root ${DIR_CONF}${FILE_CONF} && /bin/chmod 0440 ${DIR_CONF}${FILE_CONF}

endef

define fConfFiles
${DIR_CONF}${FILE_CONF}

endef

define fPreRm

endef

define fPostRm

endef