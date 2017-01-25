#!/usr/bin/env bash

if [ "$TRAVIS_BRANCH" != "master" ]; then
  echo "This commit was made against the $TRAVIS_BRANCH and not the master! No deploy!"
  exit 0
fi

if [ "$TRAVIS_PULL_REQUEST" != false ]; then
  exit 0
fi

# Is this not a build which was triggered by setting a new tag?
if [ -z "$TRAVIS_TAG" ]; then
  echo -e "Starting to tag commit.\n"

  git config --global user.email "support@e154.ru"
  git config --global user.name "delta54"

  # Add tag and push to master.
  git tag -a v${TRAVIS_BUILD_NUMBER} -m "Travis build $TRAVIS_BUILD_NUMBER pushed a tag."
  git push origin --tags
  git fetch origin

  echo -e "Done magic with tags.\n"
fi


