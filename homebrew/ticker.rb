class Ticker < Formula
  desc "Autonomous AI agent loop runner for completing epics with Claude Code"
  homepage "https://github.com/mkelk/ticker-melk"
  version "0.3.13"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.13/ticker_0.3.13_darwin_amd64.tar.gz"
      sha256 "c9d8390bd0a631cc4db36d212ce6221ff070372a3e45968f88e6291e3fddda4b"
    end
    on_arm do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.13/ticker_0.3.13_darwin_arm64.tar.gz"
      sha256 "690ec159116b5d50d972716d17aec6487756b8d03de966c8e834aae80668dfaf"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.13/ticker_0.3.13_linux_amd64.tar.gz"
      sha256 "07a7a06ae56bf1e82a929bf4670383f9eb904ce8b9bc010be2ea80a2d7427e5a"
    end
    on_arm do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.13/ticker_0.3.13_linux_arm64.tar.gz"
      sha256 "9f65231292e4e408fedd81ef0fbb7b11be2e42b5b6e6a84c12c4eb37c01d1472"
    end
  end

  def install
    bin.install "ticker"
  end

  test do
    system "#{bin}/ticker", "--version"
  end
end
