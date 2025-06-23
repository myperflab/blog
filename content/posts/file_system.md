---
title: "File System Basics"
thumbnailImagePosition: left
thumbnailImage: /images/aircraft_carrier.jpg
date: 2024-04-10
categories:
- perf engg
- perf testing
tags:
- Go Lang
showPagination: true
draft: True
---
<!--more-->

# Understanding File Systems

> **Note:** This guide focuses on the ext4 file system, one of the most widely used file systems in Linux.

## Core Concepts

Here are the fundamental concepts you need to understand about file systems:

### 1. Inode (Linux)
- Contains metadata of files
- Stores information such as:
  - File size
  - File owner
  - Block pointers
  - Access permissions
  - Timestamps

```bash
$ stat README.md 
  File: README.md
  Size: 6         	Blocks: 8          IO Block: 4096   regular file
Device: 252,1	Inode: 53767487    Links: 1
Access: (0664/-rw-rw-r--)  Uid: ( 1000/shubham-sachan)   Gid: ( 1000/shubham-sachan)
Access: 2025-06-18 11:50:34.292152078 +0530
Modify: 2025-05-10 12:59:39.211976407 +0530
Change: 2025-05-10 12:59:39.211976407 +0530
 Birth: 2025-05-10 12:59:39.211976407 +0530


```

### 2. Storage Units
- **Cluster/Block**: The smallest unit of storage allocation
  - Typically 4 KB in size
  - Helps optimize disk space usage

### 3. File System Operations
- **Mounting**: The process of making a file system accessible through a specific directory
- **Journaling**: A crash-safety feature that logs changes before committing them to disk
  - Helps maintain file system integrity
  - Prevents data corruption during system crashes

### 4. Special File Types
- **Sparse Files**: Files that efficiently store sequences of zeros by not allocating actual disk blocks
- **File References**:
  - Symbolic Links: Pointer-based references to other files
  - Hard Links: Direct references to file inodes

## File System Architecture

### 1. Basic Structure
- **Boot Block**: Contains information for bootstrapping the operating system
- **Superblock**: Stores critical metadata about the file system
- **Inode Table**: Contains metadata for all files (Unix-like systems)
- **Data Blocks**: Stores the actual contents of files

### 2. Directories and Paths
- **Path Types**:
  - Absolute Path: Starts from root (`/`)
  - Relative Path: References from current location
- **Special Directories**:
  - `.`: Current directory
  - `..`: Parent directory
  - `/dev`: Device files
  - `/proc`: Process information

### 3. File Types and Permissions

#### File Types
- Regular Files (`-`): Standard data files
- Directories (`d`): Contains file entries
- Symbolic Links (`l`): Soft links to other files
- Block Devices (`b`): Block-oriented device files
- Character Devices (`c`): Character-oriented device files
- Special Files:
  - Sockets: Inter-process communication
  - FIFO: Named pipes

#### Permission System
- **Access Rights** (`rwx`):
  - Read, Write, Execute permissions for:
    - User (owner)
    - Group
    - Others
- **Permission Management**:
  - `chmod`: Change file permissions
  - `chown`: Change file ownership
  - `umask`: Set default permissions
- **Special Permissions**:
  - SUID (Set User ID)
  - SGID (Set Group ID)
  - Sticky Bit

### 4. File Metadata
Each file maintains the following information:
- File size in bytes
- Timestamps:
  - Creation time
  - Last modification time
  - Last access time
- Ownership:
  - User owner
  - Group owner
- Permission settings
- Inode number (unique file identifier in Linux)


