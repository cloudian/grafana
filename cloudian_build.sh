#!/bin/sh -e

# Log in to quay.io
if [ -z $QUAY_IO_CREDENTIALS ]; then
  echo Must set QUAY_IO_CREDENTIALS before running this script
  exit 1
fi
USERNAME=$(echo $QUAY_IO_CREDENTIALS | cut -d: -f1)
PASSWORD=$(echo $QUAY_IO_CREDENTIALS | cut -d: -f2)
if echo "$QUAY_IO_CREDENTIALS" | grep -vq ":"; then
  echo QUAY_IO_CREDENTIALS must be in the format username:password
  exit 1
fi
echo $PASSWORD | docker login quay.io -u $USERNAME --password-stdin

cd $(dirname $0)

# Determine TAG for this build
REPO=quay.io/cloudian/grafana
TAG=$(git describe)

# Strip off Git commit for release build
if [ "$1" = "release" ]; then
  TAG=$(echo $TAG | cut -d- -f1-2)
fi

# Call standard build script.  Use --fast so it builds only for x64
packaging/docker/build.sh --fast $TAG $REPO

# Don't include leading 'v' from the repository tag
if echo "$TAG" | grep -q "^v"; then
	TAG=$(echo "$TAG" | cut -b2-)
fi

# Send to quay.io
docker push $REPO:$TAG

# If we're doing a release build, update the latest tag
if [ "$1" = "release" ]; then
  docker tag $REPO:$TAG $REPO:latest
  docker push $REPO:latest
fi