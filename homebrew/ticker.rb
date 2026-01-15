class Ticker < Formula
  desc "Autonomous AI agent loop runner for completing epics with Claude Code"
  homepage "https://github.com/mkelk/ticker-melk"
  version "0.3.12"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.12/ticker_0.3.12_darwin_amd64.tar.gz"
      sha256 "e08947a16ed8470e6a4671e232bbc3442f6c38986bc586fa21b366ac94114b4f"
    end
    on_arm do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.12/ticker_0.3.12_darwin_arm64.tar.gz"
      sha256 "5a1aaeb2950db00c24f39ab5133d563138c044cf5559ce9306e43b5b594467c9"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.12/ticker_0.3.12_linux_amd64.tar.gz"
      sha256 "21a0c93ce33e9e8dab5545af41f6eeb242157ce8acc4ecc0f64530e5803a596b"
    end
    on_arm do
      url "https://github.com/mkelk/ticker-melk/releases/download/v0.3.12/ticker_0.3.12_linux_arm64.tar.gz"
      sha256 "5ee232129a6e0e6f605dab52b5c95b6f104e7f26ac007e3b25fe28647fc46916"
    end
  end

  def install
    bin.install "ticker"
  end

  test do
    system "#{bin}/ticker", "--version"
  end
end
