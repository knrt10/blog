---
title: Seccomp security profiles for Docker
date: 2021-01-30
draft: false
summary: Understand Seccomp profiles for Docker while experimenting and creating our own seccomp profile json file and running docker container with it.
draft: false
authors:
- admin
tags:
- docker
- security
categories:
- docker
- tech
---

{{% toc %}}

## What is Seccomp?

As per the [kernel documentation](https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt)

> A large number of system calls are exposed to every userland process
with many of them going unused for the entire lifetime of the process.
As system calls change and mature, bugs are found and eradicated.  A
certain subset of userland applications benefit by having a reduced set
of available system calls.  The resulting set reduces the total kernel
surface exposed to the application.  System call filtering is meant for
use with those applications.

> Seccomp filtering provides a means for a process to specify a filter for
incoming system calls.  The filter is expressed as a Berkeley Packet
Filter (BPF) program, as with socket filters, except that the data
operated on is related to the system call being made: system call
number and the system call arguments.  This allows for expressive
filtering of system calls using a filter program language with a long
history of being exposed to userland and a straightforward data set.

## How docker uses Seccomp?

Secure computing mode `(seccomp)` is a Linux kernel feature. You can use it to restrict the actions available within the container. The `seccomp()` system call operates on the seccomp state of the calling process. You can use this feature to restrict your application’s access.

This feature is available only if Docker has been built with seccomp and the kernel is configured with `CONFIG_SECCOMP` enabled. To check if your kernel supports seccomp:

```bash
$ grep CONFIG_SECCOMP= /boot/config-$(uname -r)
CONFIG_SECCOMP=y
```

The default `seccomp` profile provides a sane default for running containers with seccomp and disables around `44` system calls out of `300+`. It is moderately protective while providing wide application compatibility. The default Docker profile can be found [here](https://github.com/moby/moby/blob/master/profiles/seccomp/default.json).

`seccomp` is instrumental for running Docker containers with least privilege. It is not recommended to change the default seccomp profile.

When you run a container, it uses the default profile unless you override it with the `--security-opt` option. For example, the following explicitly specifies a policy:

```bash
docker run --rm \
  -it \
  --security-opt seccomp=/path/to/seccomp/profile.json \
  hello-world
```

Let's take a look at snippet of syscalls allowed from the default profile:

```json
{
  "names": [
    "bpf",
    "clone",
    "fanotify_init",
    "lookup_dcookie",
    "mount",
    "name_to_handle_at",
    "perf_event_open",
    "quotactl",
    "setdomainname",
    "sethostname",
    "setns",
    "syslog",
    "umount",
    "umount2",
    "unshare"
  ],
  "action": "SCMP_ACT_ALLOW",
  "args": [],
  "comment": "",
  "includes": {
    "caps": [
        "CAP_SYS_ADMIN"
    ]
  },
  "excludes": {}
}
```

`names` field in above json snippet refers to syscalls of linux kernel. They are only allowed for containers that you run with capability `CAP_SYS_ADMIN` mentioned in `action` json field. You can pass this capability to a container using `--cap-add` flag.

## Creating out own seccomp profile json file

One thing to know is that every executable binary in unix system has some capabilities assigned to it. For example if you want to find capabilities assined to `ping` binary just use `getcap` command like this:

```bash
$ getcap $(which ping)
/usr/bin/ping = cap_net_admin,cap_net_raw+p
```

So what we will do that, is we will be tying `CAP_AUDIT_CONTROL` capability to our `chown` syscall. You can take any other capability other than `CAP_CHOWN`. 

> This is purely for experimenting and understanding how seccomp profile will work. DO NOT use it in production environment.

To make this work we will remove all the occurances of `chown` syscall from [default](https://github.com/moby/moby/blob/master/profiles/seccomp/default.json) seccomp profile and move it to our custom profile like this:

```bash
{
  "names": [
    "chown",
    "chown32",
    "fchown",
    "fchown32",
    "fchownat",
    "lchown",
    "lchown32"
  ],
  "action": "SCMP_ACT_ALLOW",
  "args": [],
  "comment": "",
  "includes": {
    "caps": [
      "CAP_AUDIT_CONTROL"
    ]
  },
  "excludes": {}
}
```

Final `json` file can be [found here](https://gist.github.com/knrt10/4c877585947f34fbfcab7626f3e7229c). Copy it from github gist and save it as `custom-profile.json` because it will be used in our next step for running docker container. Run the below command:

```bash
$ docker run --rm -it --security-opt seccomp=custom-profile.json debian bash

# Try creating a user
root@429a518f8ec5:/# useradd knrt10
useradd: failure while writing changes to /etc/passwd
```

Above command will fail as `useradd` syscall uses `CAP_CHOWN` internally. That is different topic, I will write an article about it another time.

Now exit and try to run the docker container using this command:

```bash
$ docker run --cap-add=CAP_AUDIT_CONTROL --rm -it --security-opt seccomp=custom-profile.json debian bash

# create a new user
root@ea95510fcb7c:/# useradd -m knrt10

# check current user i.e root
root@ea95510fcb7c:/# id -u
0

# create a file
root@ea95510fcb7c:/# touch a

# check permissions on file. Note currently it is owned by root
root@ea95510fcb7c:/# ls -l a
-rw-r--r-- 1 root root 0 Jan 30 11:21 a

# change ownership to earlier created user
root@ea95510fcb7c:/# chown knrt10 a

# check permissions on file. It is changed to user knrt10
root@ea95510fcb7c:/# ls -l a
-rw-r--r-- 1 knrt10 root 0 Jan 30 11:21 a
```

When running the docker container by explicitly specifying the capability `CAP_AUDIT_CONTROL` and then the container allows and uses syscall `chown`. In this way you can create your own profile and tie it up with any capability you want.

## Conclusion

You learnt how seccomp profiles are used by docker and how you can create a custom seccomp profile and use it while running your docker container. This is mostly used for security purposes when you don't want your container to have extra kernel priviledges. You can learn more in details from [docker](https://docs.docker.com/engine/security/seccomp/) documentation.

## Did you find this page helpful? Consider sharing it 🙌
