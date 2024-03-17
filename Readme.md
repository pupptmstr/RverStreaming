# Requirements Before Running the Application

To ensure the smooth operation of the application, certain prerequisites need to be met. Please follow the instructions below to set up your environment accordingly.

## Go Installation

The application is developed in Go, hence Go must be installed on your system. The version required is Go 1.15 or later. Here are brief instructions on how to install Go on various operating systems:

### Linux:

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

### macOS

For macOS, you can use Homebrew to install Go:

```bash
brew install go
```

### Windows

Download the Go installer from the [official Go website](https://go.dev/doc/install) and follow the installation instructions.


## FFmpeg Installation

This application requires FFmpeg for processing video and audio files. Here's how to install FFmpeg on different operating systems:

### Linux (Ubuntu\Debian)

Use apt to install FFmpeg:

```bash
sudo apt update
sudo apt install ffmpeg
```

### macOS

You can install FFmpeg on macOS using Homebrew:
```bash
brew install ffmpeg
```

### Windows

1) Download the FFmpeg build from [FFmpeg's official website](https://ffmpeg.org/download.html).
2) Unzip the downloaded file.
3) Add the path to the FFmpeg bin folder (e.g., C:\ffmpeg\bin) to your `PATH` Environment Variable.

## Verifying the Installation

After installing both Go and FFmpeg, you can verify the installations by running:

```bash
go version
ffmpeg -version
```

These commands should output the versions of Go and FFmpeg installed on your system, respectively. If you encounter any errors, please ensure you have correctly followed the installation instructions for your operating system.

Once all the requirements are satisfied, your system is ready to run the application.