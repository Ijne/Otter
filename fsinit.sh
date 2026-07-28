#!/bin/bash

ROOTFS="${1:-/tmp/rootfs}"
mkdir -p $ROOTFS/bin $ROOTFS/proc $ROOTFS/oldroot $ROOTFS/dev $ROOTFS/tmp
wget https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox -O $ROOTFS/bin/busybox
chmod +x $ROOTFS/bin/busybox

cd $ROOTFS/bin
ln -s busybox sh
ln -s busybox ls
ln -s busybox ps
ln -s busybox mount
ln -s busybox cat
ln -s busybox id
ln -s busybox echo
ln -s busybox clear
ln -s busybox ipcs
ln -s busybox ip


