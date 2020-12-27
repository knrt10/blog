---
title: Passing CKA exam
date: 2020-12-25
draft: false
summary: How I passed my Certified Kubenernets Administrator(CKA) exam.
draft: false
authors:
- admin
tags:
- kubernetes
- cka
- docker
- security
categories:
- kubernetes
- tech
image:
  placement: 2
  caption: "[My certification](https://www.youracclaim.com/badges/e93b0ba2-51ea-41b3-82ba-10cf4d67ac70)"
  focal_point: "Center"
  preview_only: false
  alt_text: Kautilya Tripathi's CKA certificate.
---

{{% toc %}}

## How it started

After doing a good amount of internships I started my first job working as **Cloud Infrastructure Engineer** at [Kinvolk](https://kinvolk.io). Initially when I started the job I had little experience with Kubernetes, I had only worked on setting up cluster on [GKE](https://cloud.google.com/kubernetes-engine) or worked on pre-setup clusters like [minikube](https://minikube.sigs.k8s.io/). When I started the job I worked on our product [Lokomotive](https://github.com/kinvolk/lokomotive) which is an open source Kubernetes distribution that ships pure upstream Kubernetes. It focuses on being minimal, easy to use, and secure by default. While working on that I got into deeper details of Kubernetes Infrastructure. Since then I have been trying to learn as much as I can into greater detail.

## How I practiced

> I did not took any course or followed any professional training. I mostly prefer going though documentation and practicing it. Rest it's up to you how you proceed. My day to day work is regarging k8s only but I did practice for 2 weeks for CKA exam to solve questions in time limit.

I have a **Mackbook air 2017 model with 8gb RAM**. So why am I telling this? 

It's because my laptop has gone to dust, it was hard for me to practice on my laptop. My situation currently is if I start a minikube cluster and open vscode on side, it takes some 20-30 seconds for just a file to save. Initially when my laptop somehow did not lag, I used to spun up Vagrant ubuntu machine to practice but that idea went to drain once I laptop started lagging more. To fix this problem what I did was that I used to spun up `t1.small.x86` machine on [Equinix Metal](https://console.equinix.com/)(Formerly Packet) and practice there.

## Explaining the whole process

- You first buy the exam anytime from [official CNCF website](https://www.cncf.io/certification/cka/). After that you have 1 year to take the exam before it expires
- You then schedule the exam within 1 year of time frame after you bought it
- You then give the exam online while been proctored. Before beginning of the exam the proctor checks all the criteria is met as per rules and after that only exam is started
- You complete the exam
- You get result after 36-48 hours of exam

## What's in the exam?

The exam consists of **15-20** questions, you are supposed to complete these questions in **2 hours** and you have to score **66%** or more in order to pass the exam. You are allowed to open only two browser tabs, one will be the exam interface and in the other tab you can open any of the allowed web pages, one of them is [k8s official docs](http://docs.kubernetes.io/). You can get a list of all the link that you can visit, in the candidates handbook and it says you can access https://kubernetes.io/docs/, https://github.com/kubernetes/ and https://kubernetes.io/blog

Candidate’s handbook says that you can copy 1–2 lines of data from the official documentation and that’s sufficient, you wont actually need to copy entire yaml file into your test. Its better to generate that yaml file by yourself and then edit that yaml file to have specific details.

As most of you are aware the CKA exam requires a considerable amount of preparation since it focuses on practical/hands-on questions/scenarios rather than a set of Mutiple choice questions.

## Basic Understanding of systemd

There are chances that some of the component of the cluster will be running as linux services and not as k8s pods, and you might have to debug those components to check any issues in them. So it will always be a plus to have good understanding of how to change configuration of, check logs(**journalctl**) of, start or stop a [systemd](https://www.freedesktop.org/wiki/Software/systemd/) service. You can easily learn from their official documentation.

## Networking

One of the most important thing in any distributed systems is Networking. I went little deep into it and learnt and practiced about:

- Switch routing(interfaces, routes, gateways)
- DNS
- Network Namespaces
- Docker Networking
- CNI
- Cluster Networking
- Pod and Service Networking


You should also be aware of some linux networking commands for example **nslookup, ping, curl or dig** in order to check the connectivity between the hosts or components.

This becomes very handy, if you have questions where you have to check the connectivity between the services or pods. 

## Linux Commands

You will be working on linux based machines so it’s always beneficial to know some basic command of linux based operating systems. For example how to redirect output to a file (>), filter somethings from a file (grep), find, get last or first rows from the output and cut command. Practice `awk` command and you should also know how processes work.

## k8s Imperative Commands

When you are asked to, let’s say create a pod, it’s not very ideal to write the entire manifests manually because its time consuming and you are likely to make some indention mistakes while writing the yaml. So its better to have the basic manifest of the resource generated and then edit that manifest with what is actually required by the question. More information can be found in [k8s docs](https://kubernetes.io/docs/concepts/overview/working-with-objects/object-management/).

## k8s Components

- The other part is understanding Kubernetes components and being able to fix and investigate clusters. Understand this: https://kubernetes.io/docs/tasks/debug-application-cluster/debug-cluster

- When you have to fix a component (like kubelet) in one cluster, just check how its setup on another node in the same or even another cluster. You can copy config files over etc
- I suggest you do [Kubernetes The Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way) once. It's not necessary to do it more often or know it all by heart, the CKA is not that complex. But KTHW helps understanding the concepts

- It can help if you install your own cluster using kubeadm in a VM or using a cloud provider and investigate the components. **I did it myself in a vagrant setup**
- Know how to use kubeadm to for example add nodes to a cluster
- Know how to create an Ingress resources
- Know how to snapshot/restore ETCD from another machine
- Know advanced scheduling: https://kubernetes.io/docs/concepts/scheduling/kube-scheduler

## Practice

Practice is key to everything, the exam is totally based on questions that you will have to perform on live cluster. So until and unless you have actually done those things on live cluster before it will be hard for you but if you have done some practice/hands on you don’t have to worry much about this.

The best way to practice k8s concepts is to get minikube installed and go through the [tasks](https://kubernetes.io/docs/tasks/) page of the official k8s docs. Once you have gone through the official k8s docs tasks you will be confident enough to appear in the exam.

There is no easy way out with the preparation practice, practice and sheer practice is the key to success here. Also a very economical usage of time and a personal strategy on how to target each question would greatly help crack this exam.

## Browser Terminal Setup

It should be considered to spend ~1 minute in the beginning to setup your terminal. In the real exam the vast majority of questions will be done from the main terminal. For few you might need to ssh into another machine. Just be aware that configurations to your shell will not be transferred in this case.

### Minimal Setup

#### Alias

You might have read in most of the articles to set up aliases for the exam. But personally I find it useless. I would suggest to minimally setup this alias:

```bash
echo 'alias kc="kubectl"' >> ~/.bashrc && source ~/.bashrc && echo "source <(kubectl completion bash)" >> ~/.bashrc && complete -F __start_kubectl kc
```

#### Vim

Be great with vim.
##### toggle vim line numbers

When in *vim* you can press **Esc** and type `:set number` or `:set nonumber` followed by **Enter** to toggle line numbers. This can be useful when finding syntax errors based on line - but can be bad when wanting to mark&copy by mouse. You can also just jump to a line number with `Esc :22 + Enter`.

##### Indent multiple lines

To indent multiple lines press Esc and type `:set shiftwidth=2`. First mark multiple lines using `Shift v` and the `up/down` keys. Then to indent the marked lines press `> or <`. You can then press `.` to repeat the action.

Also optionally you can create the file `~/.vimrc` with the following content:

```bash
set tabstop=2
set expandtab
set shiftwidth=2
```

The *expandtab* make sure to use spaces for tabs. Memorize these and just type them down. You can't have any written notes with commands on your desktop etc.

## Be fast

Use the `history` command to reuse already entered commands or use even faster history search through `Ctrl r` .

If a command takes some time to execute, like sometimes `kubectl delete pod x`. You can put a task in the background using `Ctrl z` and pull it back into foreground running command `fg` or run `bg` to run in background while you perform other tasks.

You can delete pods fast with:

```bash
kc delete pod x --grace-period 0 --force
```

## CKA Preparation

### Read the Curriculum

https://github.com/cncf/curriculum

### Read the Handbook

https://docs.linuxfoundation.org/tc-docs/certification/lf-candidate-handbook

### Read the important tips

https://docs.linuxfoundation.org/tc-docs/certification/tips-cka-and-ckad

### Read the FAQ

https://docs.linuxfoundation.org/tc-docs/certification/faq-cka-ckad


## Conclusion

I would say it's pretty easy if you have practiced and have good understanding how k8s works and how it is setup. It would help you while debugging and seting up cluster. Try **solving questions with most weightage first**. Relax and do not feel pressure about the time. Also you should keep in mind that **`CKA is just a tag`**. Kubernetes is very vast, so you should try to learn it as much as you can if you are interested in working in distributed systems domain.

## Did you find this page helpful? Consider sharing it 🙌
