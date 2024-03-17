# RverStreaming

Welcome to RverStreaming. This educational pet project is designed to explore the intricate world of media streaming, focusing particularly on the adaptive streaming technique. Adaptive streaming is a method that dynamically adjusts the quality of a video delivered to a viewer according to changing internet speeds and device capabilities, ensuring an optimal viewing experience.

The goal of this project is to build a practical understanding of how adaptive streaming works and to implement a simple streaming server capable of serving video and audio content in various qualities. By diving into this project, I aim to demystify the processes behind media streaming platforms and enhance my skills in developing scalable and efficient media delivery applications.

## Project Overview

The core of this project revolves around the use of Go and FFmpeg, leveraging Go's concurrency features and FFmpeg's powerful media processing capabilities to create a server that can transcode, segment, and stream media files adaptively. Whether you're interested in the technical details of setting up an adaptive streaming server, or you're looking to understand the challenges and solutions in streaming media content, this project serves as a hands-on guide to the fascinating world of adaptive streaming.

Through this project, you'll get insights into:

- The architecture and components of an adaptive streaming platform.
- How to transcode media files into various resolutions and bitrates.
- Creating DASH (Dynamic Adaptive Streaming over HTTP) manifests to describe media presentations.
- Serving media segments and manifests efficiently to clients.

Now, let's make sure your environment is set up correctly to get the most out of this project.

## Requirements Before Running the Application

To ensure the smooth operation of the application, certain prerequisites need to be met. Please follow the instructions below to set up your environment accordingly.

### Go Installation

The application is developed in Go, hence Go must be installed on your system. The version required is Go 1.15 or later. Here are brief instructions on how to install Go on various operating systems:

#### Linux:

You can install Go on Linux by using the following commands:

```bash
wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
```

Then, set up your Go environment by adding these lines to your ~/.profile or ~/.bashrc file:

```bash
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
```

#### macOS

For macOS, you can use Homebrew to install Go:

```bash
brew install go
```

#### Windows

Download the Go installer from the [official Go website](https://go.dev/doc/install) and follow the installation instructions.


### FFmpeg Installation

This application requires FFmpeg for processing video and audio files. Here's how to install FFmpeg on different operating systems:

#### Linux (Ubuntu\Debian)

Use apt to install FFmpeg:

```bash
sudo apt update
sudo apt install ffmpeg
```

#### macOS

You can install FFmpeg on macOS using Homebrew:
```bash
brew install ffmpeg
```

#### Windows

1) Download the FFmpeg build from [FFmpeg's official website](https://ffmpeg.org/download.html).
2) Unzip the downloaded file.
3) Add the path to the FFmpeg bin folder (e.g., C:\ffmpeg\bin) to your `PATH` Environment Variable.

### Verifying the Installation

After installing both Go and FFmpeg, you can verify the installations by running:

```bash
go version
ffmpeg -version
```

These commands should output the versions of Go and FFmpeg installed on your system, respectively. If you encounter any errors, please ensure you have correctly followed the installation instructions for your operating system.

Once all the requirements are satisfied, your system is ready to run the application.