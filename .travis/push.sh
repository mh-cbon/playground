#!/bin/sh

generate_files() {
  mkdir ${OUT}
  go run screenshot.go screenshot_notwindows.go -browser ${BROWSER} -action initial -out ${OUT}/initial.png
  go run screenshot.go screenshot_notwindows.go -browser ${BROWSER} -action checkbox -out ${OUT}/checkbox.png
}

setup_git() {
  git config --global user.email "travis@travis-ci.org"
  git config --global user.name "Travis CI"
}

commit_files() {
  git checkout -b gh-pages
  git add . ${OUT}/*
  git commit --message "Travis build: ${TRAVIS_BUILD_NUMBER}"
}

upload_files() {
  git remote add origin-pages https://${GH_TOKEN}@github.com/${TRAVIS_REPO_SLUG}.git > /dev/null 2>&1
  git push --quiet --set-upstream origin-pages gh-pages
}

generate_files
setup_git
commit_files
upload_files
