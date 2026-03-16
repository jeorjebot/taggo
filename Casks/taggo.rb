cask "taggo" do
  arch arm: "arm64", intel: "amd64"

  version "1.1.0"
  sha256 arm:   "114b97d97ac39d2e068792e0e03f7004259427d87004b1cc76b1bde55b5adc7d",
         intel: "ff1be1173cc42989ccc81bf70918670e71fa9b13c2e511f4d1e1e215d566e71c"

  url "https://github.com/jeorjebot/taggo/releases/download/v#{version}/taggo-v#{version}-macos-#{arch}"

  name "Taggo"
  desc "Easy peasy git tag utility"
  homepage "https://github.com/jeorjebot/taggo"

  binary "taggo-v#{version}-macos-#{arch}", target: "taggo"
end
