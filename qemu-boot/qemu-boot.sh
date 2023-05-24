#!/bin/bash -e

VERSION="6.3.3" # might be better as an argument to the script, depending on script's use case
DIRECTORY="linux-${VERSION}"
FILENAME="${DIRECTORY}.tar"
TARBALL_URL="https://cdn.kernel.org/pub/linux/kernel/v6.x/${FILENAME}.xz" 
BUSYBOX="busybox-1.36.1"
BUSYBOX_URL="https://busybox.net/downloads/${BUSYBOX}.tar.bz2"
PREREQUISITES="xz-utils curl build-essential libncurses-dev bison flex libssl-dev libelf-dev qemu busybox-static"
PROJECT_ROOT=`pwd`

install_prerequisites () {
  echo "Installing prerequisites ... "
  sudo apt install $PREREQUISITES
  echo " --> done."
}

download_and_unpack () {
  echo -n "Downloading tarball $TARBALL_URL ... "
  if [[ -f $FILENAME || -f "${FILENAME}.xz" ]]; then
    echo "tarball already downloaded--skipping."
  else
    curl -s -o ${FILENAME}.xz $TARBALL_URL
    echo "done."
  fi

  if [[ -f "${FILENAME}.xz" ]]; then
    echo -n "Decompressing tarball ... "
    unxz "${FILENAME}.xz"
    echo "done."
  fi

  if [[ -d "${DIRECTORY}" ]]; then
    echo "Found existing ${DIRECTORY} directory, skipping unpack step."
  else
    echo -n "Unpacking $FILENAME ... "
    tar xf $FILENAME
    echo "done."
  fi

  # If we wanted to verify signatures of files we've downloaded, we could do that here also. 
  # I have skipped it in the interests of time. Here is a rough outline of steps
  # curl -s -o ${FILENAME}.sign $PGP_SIGNATURE_URL # would make this a similar variable to TARBALL_URL
  # gpg --verify ${FILENAME}.sign # check the exit status of this

  echo "Downloading Busybox"
  if [[ -d ${BUSYBOX} ]]; then
    echo " ... found existing directory ${BUSYBOX}, skipping download."
  else
    curl -s -o ${BUSYBOX}.tar.bz2 $BUSYBOX_URL
    tar xjf ${BUSYBOX}.tar.bz2
    echo " --> done."
  fi
}

build_kernel () {
  echo -n "Copying kernel config ... "
  cp kernel.config ${DIRECTORY}/.config
  echo "done."

  echo "Building kernel ... "
  cd $PROJECT_ROOT/$DIRECTORY
  make -j $(nproc)
  cp arch/x86/boot/bzImage $PROJECT_ROOT
  echo " --> done."
}

build_filesystem () {
  echo "Building busybox ... "
  cp $PROJECT_ROOT/busybox.config $PROJECT_ROOT/$BUSYBOX/.config
  cd $PROJECT_ROOT/$BUSYBOX
  make -j $(nproc)

  # busybox wants to be setuid root
  sudo chown root busybox
  sudo chmod u+s busybox
  echo " --> done."

  echo -n "Configuring busybox and making filesystem ... "
  make install
  cd _install
  mkdir -p dev proc sys

  # Init is how we say hello world
  cp $PROJECT_ROOT/init .
  chmod +x init
  find . -print0 | cpio --null -ov --format=newc | gzip -9 > ${PROJECT_ROOT}/initramfs.cpio.gz
  echo "done."
}

run_qemu () {
  echo "Running QEMU"
  cd $PROJECT_ROOT
  qemu-system-x86_64 -kernel bzImage -initrd initramfs.cpio.gz
  echo "done."
}

install_prerequisites
download_and_unpack
build_kernel
build_filesystem
run_qemu

