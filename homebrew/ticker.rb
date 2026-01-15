class Ticker < Formula
  desc "Autonomous AI agent loop runner for completing epics with Claude Code"
  homepage "https://github.com/mkelk/ticker-melk"
  version "0.3.14"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.14/ticker_0.3.14_darwin_amd64.tar.gz"
      sha256 "e03156ea8487ead54dc66684db8a15eca393ce43a2bdfb05e037f39f9be9611a"
    end
    on_arm do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.14/ticker_0.3.14_darwin_arm64.tar.gz"
      sha256 "f97659cbe4fad95da41ac9d9c3f1536de180c243f0b31a8895f9eadec93a35d7"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.14/ticker_0.3.14_linux_amd64.tar.gz"
      sha256 "53e97e20d5aacb71ee9fe635ef8610894746415988106c05aa4f23a1dc3aedaf"
    end
    on_arm do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.14/ticker_0.3.14_linux_arm64.tar.gz"
      sha256 "780d52e36568c4915ab8c6ce864afe0421858bb8a9de150d14de5a3049aa8848"
    end
  end

  def install
    bin.install "ticker"
  end

  test do
    system "#{bin}/ticker", "--version"
  end
end
