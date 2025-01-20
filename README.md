# RustMaps CLI (Unofficial)

<a href="https://discord.gg/mainloot"><img src="./docs/images/mainloot-logo.png" alt="Mainloot Logo" style="width: 50%; height: auto;"></a> **&**
<a href="https://rustmaps.com"><img src="./docs/images/rustmaps.svg" alt="Mainloot Logo" style="width: 50%; height: auto;"></a>

![CI Status](https://github.com/maintc/rustmaps-cli/actions/workflows/build.yml/badge.svg)
[![Go Coverage](https://github.com/maintc/rustmaps-cli/wiki/coverage.svg)](https://raw.githack.com/wiki/maintc/rustmaps-cli/coverage.html)
![GitHub Release](https://img.shields.io/github/v/release/maintc/rustmaps-cli)

## Table of Contents
1. [Overview](#-overview)
2. [Features](#%EF%B8%8F-features)
3. [Installation](#-installation)
4. [Usage](#%EF%B8%8F-usage)
    - [Setting the API Key](#-setting-your-api-key)
    - [Generating maps](#%EF%B8%8F-generate-maps)
        - [Generate a procedural map with seed and size](#generate-a-procedural-map-with-seed-and-size)
        - [Generate a procedural map with seed and size (staging branch)](#generate-a-procedural-map-with-seed-and-size-staging-branch)
        - [Generate a procedural map with random seed](#generate-a-procedural-map-with-random-seed)
        - [Generate a custom map with seed and size](#generate-a-custom-map-with-seed-and-size)
        - [Generate a custom map with random seed](#generate-a-custom-map-with-random-seed)
        - [Generate maps from a csv file (procedural and custom)](#generate-maps-from-a-csv-file-procedural-and-custom)
        - [Download generated maps](#download-generated-maps)
        - [Download generated maps to a specified directory](#download-generated-maps-to-a-specified-directory)
    - [Opening maps in the browser](#-opening-maps-in-the-browser)
    - [Using a `csv` file](#-using-a-csv-file)
6. [Storage Locations](#-file-structurelocations)
7. [Disclaimers](#%EF%B8%8F-disclaimers)

## 📖 Overview

At [mainloot](https://mainloot.com), we use [RustMaps](https://rustmaps.com) to generate our custom maps. RustMaps provides an excellent service at a reasonable price, consider supporting them if you're a server owner. To get the most out of this tool like custom maps, you'll need at least a [premium](https://rustmaps.com/pricing) subscription.

If you want to generate maps on the command line in bulk, this tool may help you. We use this tool to generate a new map for every wipe on each of our servers. Visit us on [https://discord.gg/mainloot](https://discord.gg/mainloot).

## ⚙️ Features

| Feature                | Supported | Notes                                           |
|------------------------|-----------|-------------------------------------------------|
| Map Generator          | ✅        | Fully supported. Generate maps. No additional config required. |
| Custom Map Generator   | ✅        | Fully supported. Generate customized maps. Uses subscriber features of the API. |
| Download Maps          | ✅        | Fully Supported. Download maps and images locally from RustMaps. |
| GitHub Action (cron)   | 🚧        | Coming Soon, uses actions to automate map generation on schedule |

## 🔧 How it works
This tool takes map parameter input either via command line or `csv` file (columns: `seed`, `size`, and `saved_config`) and generates the corresponding map on [rustmaps.com](https://rustmaps.com), once completed the tool downloads the map files locally. This tool manages state files for each map. 

## 💻 Installation

### Quick Install

If you just want to get up and running quickly, the project provides a binary for `macOS`, `Linux`, and `Windows`. You can download the binary for your platform from the [releases](https://github.com/maintc/rustmaps-cli/releases) page.

### Developers and golang people

If you're familar with `go` or do not want to download a binary we recommend building from source. 

```sh
go build -o rustmaps ./
```

## ⌨️ Usage

```sh
RustMaps CLI

Usage:
  rustmaps [flags]
  rustmaps [command]

Available Commands:
  auth        Authenticate with RustMaps API
  completion  Generate the autocompletion script for the specified shell
  generate    Generate custom and procedural maps
  open        Open generated maps in the browser

Flags:
  -h, --help               help for rustmaps
  -l, --log-level string   Log level (debug, info, warn, error, dpanic, panic, fatal) (default "fatal")

Use "rustmaps [command] --help" for more information about a command.

  Resource             Path
  --------             ----
  Downloads directory  /Users/user/.rustmaps/downloads
  Imports directory    /Users/user/.rustmaps/imports
  Config file          /Users/user/.rustmaps/config.json
  Log file             /Users/user/.rustmaps/generator.log
```

### 🔑 Setting your API key

Set your api key, you can find yours at https://rustmaps.com/dashboard

```sh
rustmaps auth <rustmaps-api-key>
```

On successful login you will recieve the following message with your subscription status

```sh
API key verified: 🗺️ Premium Subscriber
```

### 🗺️ Generate maps

#### **Generate a procedural map with seed and size**

```sh
rustmaps generate --size 5000 --seed 2083170721
```

#### **Generate a procedural map with seed and size (staging branch)**

```sh
rustmaps generate --size 5000 --seed 2083170721 --staging
```

#### **Generate a procedural map with random seed**

```sh
rustmaps generate --size 5000 -r
```

#### **Generate a custom map with seed and size**

```sh
rustmaps generate --size 5000 --seed 2083170721 --saved-config default
```

#### **Generate a custom map with random seed**

```sh
rustmaps generate --size 5000 -r --saved-config default
```

#### **Generate maps from a csv file (procedural and custom)**

```sh
rustmaps generate --csv ./mymaps.csv
```

#### **Download generated maps**

You can specify `-d` to download maps after generating 

```sh
rustmaps generate ... -d
```

`...` supports all the same flags as `generate`

#### **Download generated maps to a specified directory**

You can specify `-o` to download maps to a specified directory

```sh
rustmaps generate ... -d -o ./mymaps
```

### 🌐 Opening maps in the browser

If a procedural map has already been generated on RustMaps you will not be able to generate it again. To verify this you can use the open command, this will open the map in the browser. `open` takes all the same map parameters as `generate`

```sh
rustmaps open -s 2083170721 -z 5000
```
If a URL for a map is not found, you will need to generate the map first.

```sh
rustmaps generate -s 2083170721 -z 5000 -S default
```

```sh
rustmaps open -s 2083170721 -z 5000 -S default 
```

Specifying a csv will open a dropdown allowing you to open many maps
```sh
rustmaps generate -c ./mymaps.csv
```

```sh
rustmaps open -c ./mymaps.csv
```

## 📚 Using a `csv` file

A `saved_config` value must be specified to generate a custom map, even the default. Rows with omitted `saved_config` are treated as a regular procedural map.

```csv
seed,size,saved_config
1986142550,4250,CombinedOutpost
1254873764,4250,default
719690435,4250
```

- The first map is a custom map using a custom configuration named "CombinedOutpost" (you can configure your own at https://rustmaps.com/dashboard/generator/custom)
- The second map is a custom map using the RustMaps default configuration named "default". This should be setup for you by default.
- The third map is a regular procedural map.

## 📁 File structure/locations

Run `rustmaps` by itself to see the actual paths (see [usage](#-usage))
```yaml
Resource             Path
--------             ----
Config file:         Where the rustmaps-cli configuration file lives (holds your api key)
Downloads directory: Where rustmaps-cli downloads maps/images after generation
Imports directory:   Where rustmaps-cli saves information on maps
Log file:            Where rustmaps-cli will write logs
```

You can override the `Downloads directory` with `-o`

```sh
rustmaps generate --size 5000 --seed 2083170721 -d -o ~/forcewipe
```

## ⚠️ Disclaimers

- Mainloot is not affiliated with Rustmaps.com, we're just users/fans
- This tool adheres to concurrent and monthly limits in addition to a 60 requests per minute ratelimit

## License

- [MIT](./license) (c) [mainloot](https://mainloot.com)
- [Contributing](.github/contributing.md)
