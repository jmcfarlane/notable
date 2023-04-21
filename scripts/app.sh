#!/bin/bash -e

cd $(dirname $0)/..

# App path and name
app_path=$1
app_name=$(basename "$app_path")

app_version=$2

# Icon path and name
icon_path=$3
icon_name=$(basename "$icon_path")

mkdir -p "$app_path.app/Contents/"{MacOS,Resources}

cat > "$app_path.app/Contents/Info.plist" <<END
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleGetInfoString</key>
  <string>$app_name</string>
  <key>CFBundleExecutable</key>
  <string>$app_name</string>
  <key>CFBundleIdentifier</key>
  <string>com.github.jmcfarlane.Notable</string>
  <key>CFBundleName</key>
  <string>$app_name</string>
  <key>CFBundleIconFile</key>
  <string>icon.icns</string>
  <key>CFBundleShortVersionString</key>
  <string>$app_version</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>IFMajorVersion</key>
  <integer>0</integer>
  <key>IFMinorVersion</key>
  <integer>1</integer>
</dict>
</plist>
END

cp $icon_path "$app_path.app/Contents/Resources/"
cd "$app_path.app/Contents/Resources/"

# Linux:
convert -resize 16x16   $icon_name icon_16x16.png
convert -resize 32x32   $icon_name icon_32x32.png
convert -resize 128x128 $icon_name icon_128x128.png
convert -resize 256x256 $icon_name icon_256x256.png
convert -resize 512x512 $icon_name icon_512x512.png
echo "> TIP: The next line requires: (sudo dnf install libicns-utils):"
png2icns icon.icns icon_*.png
rm -f icon*.png $icon_name
