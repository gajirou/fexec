FROM golang:1.20-bullseye

RUN curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o "/tmp/session-manager-plugin.deb" \
 && dpkg -i /tmp/session-manager-plugin.deb