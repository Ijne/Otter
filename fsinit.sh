#!/bin/bash

mkdir -p /tmp/rootfs/bin /tmp/rootfs/proc /tmp/rootfs/oldroot /tmp/rootfs/dev /tmp/rootfs/tmp
wget https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox -O /tmp/rootfs/bin/busybox
chmod +x /tmp/rootfs/bin/busybox

cd /tmp/rootfs/bin
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


