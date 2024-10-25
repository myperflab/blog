---
title: "LSM Tree Explained"
thumbnailImagePosition: left
thumbnailImage: /images/lsm.jpeg
date: 2024-10-24
categories:
- Data Security
tags:
- Ransomware
showPagination: true
draft: False
---
<!--more-->

LSM Tree -- Log-Structured Merge Tree

An LSM Tree is a data structure commonly used in databases to efficiently handle high-volume write operations. It’s particularly popular in systems where data is mostly written in large quantities (e.g., time-series databases, key-value stores, and log management systems). LSM Trees work by organizing data in layers and leveraging a combination of in-memory and on-disk data structures to ensure very much writes efficient without dropping read speed. However, read speed will be not as efficiet as in rdbms (e.g. postgress DB) as those are read optimized. In lsm tree implementation, read speed depends on the multiple parameters which will be discussed in the next sections.

