---
title: "Powerful cat /proc. How useful it can be for beginners in analyzing any issue or bug"
thumbnailImagePosition: left
thumbnailImage: /images/proc.jpeg
date: 2023-01-12
categories:
- performance engineering
tags:
- Linux
showPagination: true
draft: true
---
<!--more-->




cat /proc/PID/status
cat /proc/PID/oom_score
cat /proc/PID/cwd


Post performance test execution, the most common task we do is analyzing the server side performance metrics of the underlying hardware, performance statistics of the application (if we are capturing any) and client side performance metrics.

And if test results are meeting the SLA, well, shutdown the laptop and take a nap 😝

But in case SLA are not being met, we have two ways to look for :

1. Hardware performance and capacity
2. Application architecture and bottlenecks

## Option#1: Hardware performance and capacity

We usually capture the server side underlying hardware performance metrics which includes (basic ones) are CPU, Memory, Network and Disk performance stats. There are multiple tools and command line utility and tools which can capture the stats easily.

To understand the hardware utilization during the test, we look for the various parameters based on the kind of issues, but often we over look the disk performance stats. Today, I am trying zoom in the disk performance how to conclude if the disk was the bottleneck.

Lets assume, the application under test doing some IO operations on the disk, and you have following average numbers for disk (for example):

```bash
Disk Throughput-
	MB Read per sec: 10
	MB Write per sec: 6

Average disk queue length: 5

Disk Latency-
	Sec/read: 0.01 Sec 
	Sec/write: 0.05 Sec

IOPS: 2000
```

**Does the above numbers looks the bottleneck?** This is the tough\big question. Theoretically it can be, as we see big disk queue length here. But it may be not.

To avoid such confusion, we should test the disk performance separately. And based on those results, we may be able to say if the disk was slow or application is performing the IO operation either incorrectly or in unoptimized way.

Thanks to Jens Axboe who has written a wonderful utility **fio** **(Flexible IO Tester)** to test the disk performance. This utility can be used across the platforms. for more details checkout the [git](https://github.com/axboe/fio) page.

This utility let us execute various types of tests, and provide detailed results. These results includes the following

Bandwidth

IOPS

Latency

CPU consumption 

## **fio** **(Flexible IO Tester)**

We will go through the utility by running the  **fio** command and explaining the output:

Note: I am using ubuntu Linux VM. and to install this, run the following

```bash
sudo apt-get update
supd apt install fio
```

**In Brief,  about fio:** It is a tool that will spawn a number of threads or processes doing a particular type of IO action as specified by the user. fio takes a number of parameters\arguments (we will see in test sections), each inherited by the thread unless otherwise parameters given to them overriding that setting is given. The typical use of fio is to write a job file matching the IO load one wants to simulate.

Let's begin our test: 

**Test#1**

```bash
root@shubham-ubuntu:/home/ubuntu/disk2/fio_test# fio --name=formyblog --bs=4K --numjobs=10 --size=20M --rw=randwrite --directory=/home/ubuntu/disk2/fio_test  --ioengine=libaio --group_reporting --iodepth=16 --direct=1
formyblog: (g=0): rw=randwrite, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=16
...
fio-3.1
Starting 10 processes
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
Jobs: 4 (f=4): [_(2),w(1),_(1),w(2),_(3),w(1)][75.0%][r=0KiB/s,w=44.1MiB/s][r=0,w=11.3k IOPS][eta 00m:02s]
formyblog: (groupid=0, jobs=10): err= 0: pid=12832: Sun Jul 25 23:13:55 2021
  write: IOPS=10.4k, BW=40.6MiB/s (42.6MB/s)(200MiB/4926msec)
    slat (usec): min=3, max=103545, avg=144.65, stdev=1979.41
    clat (usec): min=125, max=534672, avg=13684.88, stdev=35875.24
     lat (usec): min=141, max=534680, avg=13829.90, stdev=36182.32
    clat percentiles (usec):
     |  1.00th=[   247],  5.00th=[   363], 10.00th=[   461], 20.00th=[   709],
     | 30.00th=[  1012], 40.00th=[  1549], 50.00th=[  2933], 60.00th=[  5538],
     | 70.00th=[  8160], 80.00th=[ 12911], 90.00th=[ 22152], 95.00th=[ 70779],
     | 99.00th=[185598], 99.50th=[231736], 99.90th=[341836], 99.95th=[341836],
     | 99.99th=[534774]
   bw (  KiB/s): min=  128, max=16432, per=9.80%, avg=4073.95, stdev=4261.40, samples=86
   iops        : min=   32, max= 4108, avg=1018.34, stdev=1065.36, samples=86
  lat (usec)   : 250=1.09%, 500=10.94%, 750=9.42%, 1000=8.29%
  lat (msec)   : 2=14.37%, 4=9.56%, 10=19.87%, 20=15.05%, 50=5.85%
  lat (msec)   : 100=1.37%, 250=3.87%, 500=0.30%, 750=0.03%
  cpu          : usr=0.40%, sys=4.20%, ctx=30632, majf=0, minf=105
  IO depths    : 1=0.1%, 2=0.1%, 4=0.1%, 8=0.2%, 16=99.7%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.1%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwt: total=0,51200,0, short=0,0,0, dropped=0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=16

Run status group 0 (all jobs):
  WRITE: bw=40.6MiB/s (42.6MB/s), 40.6MiB/s-40.6MiB/s (42.6MB/s-42.6MB/s), io=200MiB (210MB), run=4926-4926msec

Disk stats (read/write):
  sdb: ios=0/46893, merge=0/1320, ticks=0/602152, in_queue=602644, util=92.72%
```

 Let's understand all the input arguments\parameters and output\results line by line: 

### Executed command:

```bash
fio --name=formyblog --bs=4K --numjobs=10 --size=20M --rw=randwrite --directory=/home/ubuntu/disk2/fio_test  --ioengine=libaio --group_reporting --iodepth=16 --direct=1
```

**Command line arguments:**

```bash
--name: Test Name, its better to give name of the test. Results of the test can be save into a file

--bs: block size. 

--numjobs: number of parallel jobs. inabove example, we have 10 jobs, each job will write 
		(randomly) 20 MB (--size) in block size of 4 KB (--bs).

--size: size for which we are running the test.

--rw: to define Sequential or Random either read or write or both.
		various available option for rw:
				read:Sequential reading
				write:Sequential Write
				randread:random reading
				randwrite:random writing
				readwrite:Mixed, sequential workload
				randrw:Mixed Random Workload
		By default, for a mixed workload, the read and write percentages are set to 50% / 50%.
		In order to achieve a different distribution (%), the parameter ” –rwmixread” or ” –rwmixwrite ”
		can also be specified.

--director: path of the dir where it does read\write operations. its very important to define
		path. We can specify a number of directories by separating the names with a ‘:’ character.

--group_reporting: to publish the results grouped together. if we dont use it , 
	we will get result separatly for each job.

--iodepth: number of IOs issues to the OS at any given time. BE vary aware that, this is not 
		equivalent ot OS IO depth, and this setting is entirely application side. combining both 
		"--numjobs=x" and this parameter, fio can issues x parallel jobs and each job can issue
		"--iodepth=y" y IO at any given time. 

--ioengine: Defines how the job issues I/O. There are many, few are like:
		mmpa, sync, psync, libaio, vsync, posixaio etc.
		libaio: Linux native asynchronous I/O. Note: for async IO beahaviour "--direct=1" arg is required,
		because the page cache can not be addressed asynchronously under Linux.
	
```

There are other various arguments which can be used, few are given below. Except a very unordinary test, these arguments are more than enough for a performance engineer, but don't restricts to these only and try to explore as much as you can. Easiest way is to go through fio documentations. There are many source, few can be found in the reference section at the bottom of this page.

```bash
--zero_buffers: Initialise buffers with all zeros. Default: fill buffers with random data.
		usually fio generate random data (compressible) before the disk IO operations and keep into buffers.
		incase we want to avoid compression effect during the test, we can direct fio torefill the buffers 
		afer every IO submit with the the help of "-–refill_buffers" arguments.
	
--runtime: (deafult unit second) Tell fio to terminate processing after the specified period of time. 
		It can be quite hard to determine for how long a specified job will run, so this parameter is 
		handy to cap the total runtime to a given time. 

--time_based: If set, fio will run for the duration of the runtime specified even if the file
		are completely read or written. It will simply loop over the same workload as many times as 
		the runtime allows.

--output=filename (Write output to file filename.). this can also be achived simply putting 
			">>" at the end and giveng file name.
```

### Command Output:

 When fio test starts, it create the parallel process and send the updates of each threads, but not all the update you will see on the console. Like the below one these. logs are indicates the job is running.

```bash
formyblog: (g=0): rw=randwrite, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=16
...
fio-3.1
Starting 10 processes
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
formyblog: Laying out IO file (1 file / 20MiB)
Jobs: 4 (f=4): [_(2),w(1),_(1),w(2),_(3),w(1)][75.0%][r=0KiB/s,w=44.1MiB/s][r=0,w=11.3k IOPS][eta 00m:02s]

```

Last line of the above is the update\status from each process started above. Actual output of the fio starts after this string. 

In the last line of above, r=0 stands for read is 0, as we are doing write io operation, value of r will be zero. w stands for write.

The next set lines of the output (pasted below for easy access) is not the summary of all the jobs, **instead these lines belong to one of the multiple jobs (--numjobs=10 in our case).** If you see the "pid" it belongs to one job. But few parameters in the output are either summary or aggregated value of the all jobs. e.g. "BW" (e.g. BW=40.6MiB/s) is the aggregated value of all the jobs, but "bw" (e.g. bw (  KiB/s): min=  128...) is only for a particular job. Similarly IOPS and iops, give a thought. 

```bash
formyblog: (groupid=0, jobs=10): err= 0: pid=12832: Sun Jul 25 23:13:55 2021
  write: IOPS=10.4k, BW=40.6MiB/s (42.6MB/s)(200MiB/4926msec)
    slat (usec): min=3, max=103545, avg=144.65, stdev=1979.41
    clat (usec): min=125, max=534672, avg=13684.88, stdev=35875.24
     lat (usec): min=141, max=534680, avg=13829.90, stdev=36182.32
    clat percentiles (usec):
     |  1.00th=[   247],  5.00th=[   363], 10.00th=[   461], 20.00th=[   709],
     | 30.00th=[  1012], 40.00th=[  1549], 50.00th=[  2933], 60.00th=[  5538],
     | 70.00th=[  8160], 80.00th=[ 12911], 90.00th=[ 22152], 95.00th=[ 70779],
     | 99.00th=[185598], 99.50th=[231736], 99.90th=[341836], 99.95th=[341836],
     | 99.99th=[534774]
   bw (  KiB/s): min=  128, max=16432, per=9.80%, avg=4073.95, stdev=4261.40, samples=86
   iops        : min=   32, max= 4108, avg=1018.34, stdev=1065.36, samples=86
  lat (usec)   : 250=1.09%, 500=10.94%, 750=9.42%, 1000=8.29%
  lat (msec)   : 2=14.37%, 4=9.56%, 10=19.87%, 20=15.05%, 50=5.85%
  lat (msec)   : 100=1.37%, 250=3.87%, 500=0.30%, 750=0.03%
  cpu          : usr=0.40%, sys=4.20%, ctx=30632, majf=0, minf=105
  IO depths    : 1=0.1%, 2=0.1%, 4=0.1%, 8=0.2%, 16=99.7%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.1%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwt: total=0,51200,0, short=0,0,0, dropped=0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=16
```

**Line (1) (from the above): "formyblog: (groupid=0, jobs=10): err= 0: pid=12832: Sun Jul 25 23:13:55 2021"**

```bash
formyblog: is the test name we have given as argument.
grooupid: by deafult it is 0. 
err: Error count if occured, esle it will be zero
pid: Proicess Id of the fio test
at the enf we get the time when job starts
```

**Line (2): "write: IOPS=10.4k, BW=40.6MiB/s (42.6MB/s)(200MiB/4926msec)"**

```bash
as we have executed write operation (--rw=randwrite) we get write IOPS and Bandwidth (BW)
```

**Line (3, 4, 5 ): slat, clat and lat**

**"slat (usec): min=3, max=103545, avg=144.65, stdev=1979.41"
"clat (usec): min=125, max=534672, avg=13684.88, stdev=35875.24"
"lat (usec): min=141, max=534680, avg=13829.90, stdev=36182.32"**

```bash
slat: **Submission  latency** (min, max, avg and stdev as name suggest are minimum, maximum, average and standard 
		deviation). This is the time it took to submit the I/O. For sync I/O  this  row  is  not displayed as the 
		slat is really the completion latency (since queue/complete is one operation there).

clat: **Completion latency**. This denotes the time from submission to  completion  of  the I/O pieces. 
		For sync I/O, clat will usually be equal (or very close) to 0, as the time from submit to complete is 
		basically just CPU time (I/O has already been done, see slat explanation).

lat: **Total latency**. This denotes the time from when fio created the I/O  unit to completion of the I/O operation.
```

For a high-level disk perf understanding, **lat** is pretty much enough to make conclusion.

**Line 6: clat**

```bash
clat percentiles (usec):
     |  1.00th=[   247],  5.00th=[   363], 10.00th=[   461], 20.00th=[   709],
     | 30.00th=[  1012], 40.00th=[  1549], 50.00th=[  2933], 60.00th=[  5538],
     | 70.00th=[  8160], 80.00th=[ 12911], 90.00th=[ 22152], 95.00th=[ 70779],
     | 99.00th=[185598], 99.50th=[231736], 99.90th=[341836], 99.95th=[341836],
     | 99.99th=[534774]
```

Above shows clat on percentile basis. and show latency in usec (microseconds). As a performance engineer, we fairly understand what percentile is and how much useful it is in concluding the performance of any app.

**Line 7: bw**

```bash
bw (  KiB/s): min=  128, max=16432, per=9.80%, avg=4073.95, stdev=4261.40, samples=86
```

"bw" - Bandwidth statistics based on samples. It also includes the number of samples taken (samples) and an approximate percentage of total aggregate bandwidth this thread received in its group (per). This last value is only really useful if the threads in this group are on the same disk, since they are then competing for disk access.

<aside>
💡 The "bw" does not shows overall bandwidth (BW). It belongs to only 1job and calculated for the given samples. This value will always be lower than overall bandwidth, as there are more than one job are running in parallel, so over all bandwidth will be shared among the process running at a time.

</aside>

**Line 8: iops**

```bash
iops        : min=   32, max= 4108, avg=1018.34, stdev=1065.36, samples=86
```

iops: input output per sec. It is straightforward. But, remember, this value is also calculated on the number of samples, and will be lower than over all **IOPS**. 

**Line 9-11: lat-** lat is latency as explained in the above section. below output shows distribution with percentage. e.g. 1.09% transactions got 250 usec latency. 

```bash
	lat (usec)   : 250=1.09%, 500=10.94%, 750=9.42%, 1000=8.29%
  lat (msec)   : 2=14.37%, 4=9.56%, 10=19.87%, 20=15.05%, 50=5.85%
  lat (msec)   : 100=1.37%, 250=3.87%, 500=0.30%, 750=0.03%
```

**Line 12: cpu -** user/system CPU percentages followed by context switches then major and minor page faults. Since this test is configured to use direct IO (—direct=1 in the command), there should be very few page faults.

```bash
cpu          : usr=0.40%, sys=4.20%, ctx=30632, majf=0, minf=105
```

**Line 13-17: IO depths-** IO depth parameter basically instructs fio to submit n number of IO asynchronously and it does not have to wait for earlier I/O to be "acknowledged" before submitting new I/O until it reaches the IO depth limit. e.g. IO depth=16, fio will send 1 to 16 instruction and it does not have to wait for the response. But to send 17th, it must have received at least 1 response from the outstanding already sent instructions. 

```bash
IO depths    : 1=0.1%, 2=0.1%, 4=0.1%, 8=0.2%, 16=99.7%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.1%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwt: total=0,51200,0, short=0,0,0, dropped=0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=16
```

**submit:** How many pieces of I/O were submitting in a single submit call. Each entry denotes that amount and below, until the previous entry – e.g., 16=100% means that we submitted anywhere between 9 to 16 I/Os per submit call.

**complete:** similar to submit.

**issued rwt:** The number of read/write/trim requests issued, and how many of them were short or dropped.

There are many other IO depth parameters which can be used for emulating other scenario.

 

**Line 19-20:**  This line present the summary and final result. and it is quite self explanatory. 

```bash
Run status group 0 (all jobs):
  WRITE: bw=40.6MiB/s (42.6MB/s), 40.6MiB/s-40.6MiB/s (42.6MB/s-42.6MB/s), io=200MiB (210MB), run=4926-4926msec
```

**bw:** overall bandwidth of the disk

**io:** total size of the IO

**run:** Total time taken in test execution

At last fio also gives disk stats:

```bash
Disk stats (read/write):
  sdb: ios=0/46893, merge=0/1320, ticks=0/602152, in_queue=602644, util=92.72
```

**ios:** Number of I/Os performed by all groups.

**merge:** Number of merges performed by the I/O scheduler.

**ticks:** Number of ticks we kept the disk busy.

**in_queue:** Total time spent in the disk queue.

**util:** The disk utilization. A value of 100% means we kept the disk busy constantly, 50% would be a disk idling half of the time.

From the above exercise and output, we can see the Disk is highly efficient in writing when IO Depth is 16 and parallel jobs are up to 10. 

However, the disk perf numbers which we have assumed for an application under test show queue length is 5 and Read (MB)/sec is 10, which is still not concludable if that is the issue. 

As the test we have executed was only for write and that does not tells about how the disk will perform. We may execute another test where we do both in parallel. 

On disk queue length, There could be reason, application writing threads could be ore what we have used (10). We may want to try another test.

Next test could be like below and observe the output: 

```bash
fio --name=formyblog2 --bs=4K --numjobs=100 --size=20M --rw=readwrite --directory=/home/ubuntu/disk2/fio_test --ioengine=libaio --group_reporting --iodepth=256 --direct=1
```

I leave the output for you to do analysis. We will be seeing two set of results separately for read and write because In this test, we have done both read and write operations. Hope this gives you a head start to understand the fio output.

```bash
root@shubham-ubuntu:/home/ubuntu/disk2/fio_test# fio --name=formyblog2 --bs=4K --numjobs=100 --size=20M --rw=readwrite --directory=/home/ubuntu/disk2/fio_test --ioengine=libaio --group_reporting --iodepth=256 --direct=1
formyblog2: (g=0): rw=rw, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=256
...
fio-3.1
Starting 100 processes
Jobs: 96 (f=96): [M(3),_(1),M(45),_(1),M(21),_(1),M(4),_(1),M(23)][50.0%][r=60.5MiB/s,w=63.3MiB/s][r=15.5k,w=16.2k IOPS][eta 0Jobs: 94 (f=94): [M(1),_(1),M(1),_(1),M(8),_(1),M(36),_(1),M(21),_(1),M(4),_(1),M(23)][53.8%][r=52.6MiB/s,w=50.7MiB/s][r=13.5kJobs: 75 (f=74): [_(2),M(1),_(1),M(8),_(1),M(3),_(1),M(5),_(1),f(1),_(2),M(6),_(1),M(5),_(1),M(9),_(3),M(6),_(1),M(1),_(1),M(3),_(2),M(6),_(2),M(3),_(1),M(7),_(1),M(6),_(1),M(1),_(1),M(1),_(2),M(3)][65.2%][r=78.2MiB/s,w=78.4MiB/s][r=20.0k,w=20.1k IOPS]Jobs: 50 (f=49): [_(2),M(1),_(2),M(1),_(2),M(4),_(2),M(2),_(2),M(1),_(1),M(1),_(5),M(1),_(1),M(2),_(1),M(1),_(1),M(2),_(1),M(1),_(4),M(2),_(1),M(1),_(1),M(1),_(4),M(6),_(1),M(1),_(1),M(3),_(2),M(1),_(1),M(2),_(1),M(1),_(2),M(2),_(2),M(1),_(1),M(2),f(1),M(2),_(1),M(3),_(1),M(1),_(2),M(1),_(1),M(1),_(2),M(1),_(2)][72.7%][r=76.7MiB/s,w=76.5MiB/s][r=19.6k,w=19.6k IOPS][eta 00m:06s]
formyblog2: (groupid=0, jobs=100): err= 0: pid=25987: Sat Jul 31 17:01:39 2021
   read: IOPS=15.7k, BW=61.3MiB/s (64.3MB/s)(1002MiB/16337msec)
    slat (nsec): min=1645, max=2385.3M, avg=1989971.80, stdev=37389582.53
    clat (usec): min=228, max=3079.4k, avg=634650.80, stdev=445075.07
     lat (usec): min=244, max=3852.3k, avg=636640.92, stdev=445806.72
    clat percentiles (msec):
     |  1.00th=[    4],  5.00th=[   21], 10.00th=[   96], 20.00th=[  257],
     | 30.00th=[  330], 40.00th=[  451], 50.00th=[  625], 60.00th=[  726],
     | 70.00th=[  810], 80.00th=[  919], 90.00th=[ 1183], 95.00th=[ 1418],
     | 99.00th=[ 2165], 99.50th=[ 2433], 99.90th=[ 2735], 99.95th=[ 2937],
     | 99.99th=[ 3071]
   bw (  KiB/s): min=    6, max= 6576, per=1.69%, avg=1058.57, stdev=706.28, samples=1771
   iops        : min=    1, max= 1644, avg=264.36, stdev=176.54, samples=1771
  write: IOPS=15.6k, BW=61.1MiB/s (64.1MB/s)(998MiB/16337msec)
    slat (usec): min=2, max=2385.3k, avg=1894.46, stdev=35896.03
    clat (usec): min=545, max=3079.2k, avg=645759.96, stdev=445159.18
     lat (usec): min=548, max=3583.5k, avg=647654.57, stdev=445869.08
    clat percentiles (msec):
     |  1.00th=[    7],  5.00th=[   30], 10.00th=[  116], 20.00th=[  271],
     | 30.00th=[  347], 40.00th=[  472], 50.00th=[  634], 60.00th=[  726],
     | 70.00th=[  818], 80.00th=[  927], 90.00th=[ 1200], 95.00th=[ 1435],
     | 99.00th=[ 2232], 99.50th=[ 2433], 99.90th=[ 2802], 99.95th=[ 2937],
     | 99.99th=[ 3071]
   bw (  KiB/s): min=    5, max= 6832, per=1.71%, avg=1070.46, stdev=691.91, samples=1742
   iops        : min=    1, max= 1708, avg=267.33, stdev=172.94, samples=1742
  lat (usec)   : 250=0.01%, 500=0.01%, 750=0.03%, 1000=0.04%
  lat (msec)   : 2=0.30%, 4=0.57%, 10=0.93%, 20=2.59%, 50=2.32%
  lat (msec)   : 100=2.95%, 250=8.63%, 500=23.35%, 750=20.92%, 1000=20.04%
  lat (msec)   : 2000=15.88%, >=2000=1.46%
  cpu          : usr=0.16%, sys=0.79%, ctx=5223, majf=0, minf=1297
  IO depths    : 1=0.1%, 2=0.1%, 4=0.1%, 8=0.2%, 16=0.3%, 32=0.6%, >=64=98.8%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.1%
     issued rwt: total=256423,255577,0, short=0,0,0, dropped=0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=256

Run status group 0 (all jobs):
   READ: bw=61.3MiB/s (64.3MB/s), 61.3MiB/s-61.3MiB/s (64.3MB/s-64.3MB/s), io=1002MiB (1050MB), run=16337-16337msec
  WRITE: bw=61.1MiB/s (64.1MB/s), 61.1MiB/s-61.1MiB/s (64.1MB/s-64.1MB/s), io=998MiB (1047MB), run=16337-16337msec

Disk stats (read/write):
  sdb: ios=4684/4639, merge=251098/250336, ticks=1023060/1039920, in_queue=2071768, util=98.94%
```

**References:**

[https://fio.readthedocs.io/en/latest/index.html](https://fio.readthedocs.io/en/latest/fio_doc.html#)

[https://unix.stackexchange.com/questions/459045/what-exactly-is-iodepth-in-fio](https://unix.stackexchange.com/questions/459045/what-exactly-is-iodepth-in-fio)

[https://docs.virtuozzo.com/virtuozzo_hybrid_infrastructure_4_0_benchmarking_guide/appendix/randwrite.html](https://docs.virtuozzo.com/virtuozzo_hybrid_infrastructure_4_0_benchmarking_guide/appendix/randwrite.html)