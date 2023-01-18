---
title: "How File entropy can help in identifying ransomware attack"
thumbnailImagePosition: left
thumbnailImage: /images/s1.jpeg
date: 2023-01-10
categories:
- Data Security
tags:
- Ransomware
showPagination: true
draft: False
---
<!--more-->

Data security getting important than ever, more often ransomware attack happens and some time its severity is very high and critical. Organizations\Enterprises has no other option except paying the ransom. Recently I have stumbled on a very good blog/article about how file entropy can be used as a good measure to restrict and identify such attacks. Click [here](https://practicalsecurityanalytics.com/file-entropy/) to read more. Those who are new to data security, entropy, file signature or packers can read below to form initial understanding:

**File Entropy**: In simple words - File entropy is the measurement the randomness of the data in a file. As per the mathematical formula, entropy has a resulting value of something between zero (0) and eight (8). The closer the number is to zero, the more orderly or non-random the data is. The closer the data is to the value of eight, the more random or non-uniform the data is.

for example, if a file (1KB) contains only 0, its entropy will be 0. Same file, contains half 0 and half 1. Entropy of the file will be non-zero.

Another examples with simple python code:

```
from scipy.stats import entropy as en
import pandas as pd

sample_list = []
print(en(pd.Series(sample_list).value_counts()))
```
Using this code snippet, tried a few lists with various data type, see below

```
Entropy of list b - [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
0.0
Entropy of list c - [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1]
0.0
Entropy of list d - [0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1]
0.6931471805599453
Entropy of list e - [0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1]
0.6931471805599453
Entropy of list f - [0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0]
1.094780226007948
Entropy of list g - [0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0]
1.6020559154587264
Entropy of list h - [0, 1, 2, 3, 4, 5, 7, 8, 9, 10, 11, 12, 1, 3, 14, 15]
2.599301927099795
Entropy of list i - ['a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a']
0.0
Entropy of list j - ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j']
2.3025850929940455
Entropy of list k - ['hi', 'we', 'are', 'checking', 'entropy', 'on', 'random', 'string']
2.0794415416798357
Entropy of list l - [0, 'a', '<', '&', '9', 3, ':', 'o', '#', '*']
2.3025850929940455
```

**File Hash**: Hashes are the output of a hashing algorithm like MD5 (Message Digest 5) or SHA (Secure Hash Algorithm). These algorithms essentially aim to produce a unique, fixed-length string – the hash value, or “message digest” – for any given piece of data or “message”. As every file on a computer is, ultimately, just data that can be represented in binary form, a hashing algorithm can take that data and run a complex calculation on it and output a fixed-length string as the result of the calculation. The result is the file’s hash value or message digest. Read [here](https://www.sentinelone.com/cybersecurity-101/hashing/) more.

**Packer**: A Packer is a program that takes the executable as input, and it uses compression to obfuscate the executable's content. This obfuscated content is then stored within the structure of a new executable file; the result is a new executable file (packed program) with obfuscated content on the disk. Upon execution of the packed program, it executes a decompression routine, which extracts the original binary in memory during runtime and triggers the execution.

Why does malware use packer: Malware always wanted to go undetected in the system. To achieve this, they need something which can obfuscate the malware from antivirus scans. Packer is the solution for this. Although packer is not designed for this, but for malware it is quite required. Packer's more legitimate use is to protect intellectual property or other sensitive data from being copied.

There are many packer softwares are available in the market, but most of the malware use their own custom packer to make it difficult for the system\data engineer to retrieve the files which have been impacted due to ransomware attacks. Read [here](https://resources.infosecinstitute.com/topic/analyzing-packed-malware/#gref) more.

**PEID**: The Portable Executable (PE) file header contains the metadata about the executable file itself. At its bare minimum, it comprises of the following: a DOS stub, a signature, the architecture of the file's code, a time stamp, a pointer, and various flags.

There are many tools and ways to scan the PE header which can tells a lot about the executable\malware which can be used for the ransomware analysis.








