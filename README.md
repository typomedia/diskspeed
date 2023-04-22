# `diskspeed` - Disk speed benchmark tool

This is a fork of [gobonniego](https://github.com/cunnie/gobonniego) with some minor changes. It is a simple tool to benchmark the speed of your disk. It can be used to test the performance of a new disk or to compare the performance of different disks.

For more information about the benchmarking methodology, see [bonnie-64](https://code.google.com/archive/p/bonnie-64/) and [gobonniego](https://github.com/cunnie/gobonniego).

## Usage

    diskspeed [options]

## Options
    -d, --dir string            The directory to use for the test (default "/home/typomedia")
    -g, --gb float              The amount of disk space to use for the test (default 15)
    -h, --help                  Show help
    -i, --iops-duration float   The duration in seconds to use for the IOPS test (default 15)
    -j, --json                  Output results in JSON format
    -r, --runs int              The number of test runs (default 1)
    -s, --seconds int           The time in seconds to run the test
    -t, --threads int           The number of concurrent readers/writers. Defaults to the number of CPU cores (default 4)
    -v, --verbose               Verbose output
    -V, --version               Show version information

## Example

    ./diskspeed
    diskspeed 0.1.0
    Sequential Write MB/s: 243.59
    Sequential Read MB/s: 432.35
    IOPS: 19679

---
Copyright © 2023 Typomedia Foundation. All rights reserved.